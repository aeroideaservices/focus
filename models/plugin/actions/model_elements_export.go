package actions

import (
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	"github.com/aeroideaservices/focus/models/plugin/entity"
	"github.com/aeroideaservices/focus/models/plugin/focus"
	"github.com/aeroideaservices/focus/services/errors"
)

type Export struct {
	repository              ExportInfoRepository
	modelsRegistry          ModelsRegistry
	exporter                Exporter
	fileStorage             FileStorage
	logger                  *zap.SugaredLogger
	fileStorageBaseEndpoint string
}

func NewExport(
	repository ExportInfoRepository,
	modelsRegistry ModelsRegistry,
	exporter Exporter,
	fileStorage FileStorage,
	logger *zap.SugaredLogger,
	fileStorageBaseEndpoint string,
) *Export {
	return &Export{
		repository:              repository,
		modelsRegistry:          modelsRegistry,
		exporter:                exporter,
		fileStorage:             fileStorage,
		logger:                  logger,
		fileStorageBaseEndpoint: fileStorageBaseEndpoint,
	}
}

func (s Export) Export(ctx context.Context, action ExportModelElements) error {
	model := s.modelsRegistry.GetModel(action.ModelCode)
	if model == nil {
		return errors.NotFound.Newf("model with code \"%s\" not found", action.ModelCode)
	}

	filter := ListModelElementsQuery{
		ModelCode:    model.Code,
		SelectFields: nil,
		OrderBy:      action.OrderBy,
	}

	// преобразуем значения фильтров к нужным типам
	for fieldCode, values := range action.Filter {
		field := model.Fields.GetByCode(fieldCode)
		if field == nil || !field.Filterable {
			return errors.BadRequest.New("got wrong field code for filter")
		}

		if field.IsTime && len(values) != 2 {
			return errors.BadRequest.Newf("field \"%s\" must contains two values", fieldCode)
		}

		vs, err := field.Slice(values)
		if err != nil {
			return errors.NoType.Wrap(err, "error converting filter values")
		}

		filter.Filter.FieldsFilter[fieldCode] = vs
	}

	// получаем коды полей, которые нужно вернуть в запросе
	for _, field := range model.Fields {
		if !slices.Contains(field.Hidden, focus.ListView) || field == model.PrimaryKey { // всегда возвращаем PK
			filter.SelectFields = append(filter.SelectFields, field.Code)
		}
	}

	go s.export(ctx, model, filter)

	return nil
}

func (s Export) GetExportInfo(ctx context.Context, action GetExportInfo) (any, error) {
	s.logger.Debug("getting model export info", "dto", action)
	model := s.modelsRegistry.GetModel(action.ModelCode)
	if model == nil {
		return nil, errors.NotFound.Newf("model with code \"%s\" not found", action.ModelCode)
	}

	exportInfo, err := s.repository.GetLast(ctx, model.Code)
	if err != nil {
		return nil, err
	}

	return exportInfo, nil
}

func (s Export) export(ctx context.Context, model *focus.Model, filter ListModelElementsQuery) {
	s.logger.Debugw("start exporting model", "modelCode", model.Code)

	s.logger.Debugw("creating export info", "modelCode", model.Code)
	export := entity.ExportInfo{
		ID:        uuid.New(),
		ModelCode: filter.ModelCode,
		Status:    entity.StatusPending,
		Time:      time.Now(),
	}
	err := s.repository.Create(ctx, export)
	if err != nil {
		s.logger.Errorw("error creating export info", "err", err)
		return
	}
	s.logger.Debugw("export info has been created", "modelCode", model.Code)
	defer func() { // ставим таймер удаления записи
		go func() { defer s.delete(export.ID); time.Sleep(time.Minute * 5) }()
	}()

	s.logger.Debugw("getting model export file", "modelCode", model.Code)
	file, err := s.exporter.GetFile(ctx, model, filter)
	if err != nil {
		s.logger.Errorw("error exporting model", "err", err)
		export.Status = entity.StatusError
		err = s.repository.Update(ctx, export)
		if err != nil {
			s.logger.Errorw("error updating export info", "err", err)
		}
		return
	}
	_ = file.Close()

	fileReader, err := os.Open(file.Name())
	if err != nil {
		s.logger.Errorw("error reading file", "err", err)
		return
	}
	defer func(fileReader *os.File) { _ = fileReader.Close() }(fileReader)
	s.logger.Debugw("got temporary model export file", "modelCode", model.Code, "filename", file.Name())

	s.logger.Debugw("uploading file to file storage", "modelCode", model.Code)
	key := fmt.Sprintf("models/%s_%s%s", model.Code, time.Now().Format("2006-01-02_15:04:05"), filepath.Ext(file.Name()))
	export.Filepath = s.fileStorageBaseEndpoint + "/" + key // todo

	createFile := &CreateFile{
		Key:         key,
		ContentType: mime.TypeByExtension(filepath.Ext(file.Name())),
		File:        fileReader,
	}
	if err = s.fileStorage.Upload(ctx, createFile); err != nil {
		s.logger.Errorw("error uploading export file", "err", err)
		s.logger.Debugw("updating export info", "modelCode", model.Code)
		export.Status = entity.StatusError
		err = s.repository.Update(ctx, export)
		if err != nil {
			s.logger.Errorw("error updating export info", "err", err)
		} else {
			s.logger.Debugw("export info has been updated", "modelCode", model.Code)
		}
		return
	}
	s.logger.Debugw("file has been uploaded to file storage", "modelCode", model.Code, "filepath", export.Filepath)

	s.logger.Debugw("updating export info", "modelCode", model.Code)
	export.Status = entity.StatusSucceed
	err = s.repository.Update(ctx, export)
	if err != nil {
		s.logger.Errorw("error updating export info", "err", err)
	}
	s.logger.Debugw("export info has been updated", "modelCode", model.Code)

	s.logger.Debugw("model has been exported", "modelCode", model.Code)
}

func (s Export) delete(exportInfoId uuid.UUID) {
	s.logger.Debugw("deleting export file", "exportInfoId", exportInfoId)
	s.logger.Debugw("deleting export info", "id", exportInfoId)
	ctx := context.Background()
	exportFilepath, err := s.repository.Delete(ctx, exportInfoId)
	if err != nil {
		s.logger.Errorw("error deleting export info", "err", err)
		return
	}
	s.logger.Debugw("export info has been deleted", "id", exportInfoId)

	if exportFilepath == "" {
		s.logger.Errorw("export filepath is empty", "err", err)
		return
	}

	s.logger.Debugw("deleting export file", "filepath", exportFilepath)
	key := strings.TrimPrefix(exportFilepath, s.fileStorageBaseEndpoint+"/")
	err = s.fileStorage.Delete(ctx, key)
	if err != nil {
		s.logger.Errorw("error deleting export file", "err", err)
		return
	}
	s.logger.Debugw("export file has been deleted", "filepath", exportFilepath)
}
