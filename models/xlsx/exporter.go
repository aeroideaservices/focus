package xlsx

import (
	"context"
	"os"

	"github.com/xuri/excelize/v2"
	"golang.org/x/exp/slices"

	"github.com/aeroideaservices/focus/models/plugin/actions"
	"github.com/aeroideaservices/focus/models/plugin/focus"
	"github.com/aeroideaservices/focus/services/errors"
)

const defaultSheetName = "Sheet1"

// Exporter сервис работы с экспортом элементов модели
type Exporter struct {
	repositoryResolver actions.RepositoryResolver
	batchSize          int
}

// NewExporter конструктор
func NewExporter(repositoryResolver actions.RepositoryResolver, batchSize int) *Exporter {
	return &Exporter{
		repositoryResolver: repositoryResolver,
		batchSize:          batchSize,
	}
}

// GetFile получение файла экспорта
func (e Exporter) GetFile(ctx context.Context, model *focus.Model, filter actions.ListModelElementsQuery) (*os.File, error) {
	osFile, err := os.CreateTemp("", model.Code+"*.xlsx")
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error creating temporary file")
	}

	file := excelize.NewFile()
	defer func(file *excelize.File) { _ = file.Close() }(file)

	// переименуем дефолтный лист
	if err := file.SetSheetName(defaultSheetName, filter.ModelCode); err != nil {
		return nil, err
	}

	streamWriter, err := file.NewStreamWriter(model.Code)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error creating stream writer")
	}

	// формируем и записываем заголовки таблицы
	tableHeaders := make([]any, len(model.Fields))
	for i, field := range model.Fields {
		tableHeaders[i] = excelize.Cell{Value: field.Code}
	}
	cell, _ := excelize.CoordinatesToCellName(1, 1)
	err = streamWriter.SetRow(cell, tableHeaders)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error setting sheet row")
	}

	// пишем в файл
	if err := e.writeFile(ctx, streamWriter, model, filter); err != nil {
		return nil, err
	}
	if err := streamWriter.Flush(); err != nil {
		return nil, errors.NoType.Wrap(err, "error flushing sheet")
	}
	if err := file.Write(osFile); err != nil {
		return nil, errors.NoType.Wrap(err, "error writing file")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return osFile, nil
	}
}

func (e Exporter) writeFile(ctx context.Context, sw *excelize.StreamWriter, model *focus.Model, filter actions.ListModelElementsQuery) error {
	// получаем репозиторий модели
	repository := e.repositoryResolver.Resolve(model.Code)
	if repository == nil {
		return errors.NoType.New("cannot resolve repository")
	}
	total, err := repository.Count(ctx, filter.Filter)
	if err != nil {
		return err
	}

	// пачками получаем и записываем строки в файл
	for i := 0; i < (int(total)/e.batchSize)+1; i++ {
		filter.Offset = e.batchSize * i
		filter.Limit = e.batchSize
		var modelElements []any
		if modelElements, err = repository.List(ctx, filter); err != nil {
			return errors.NoType.Wrap(err, "error getting model list")
		}

		for j, modelElement := range modelElements {
			axis, _ := excelize.CoordinatesToCellName(1, i*e.batchSize+j+2)
			cells, err := getCells(model, modelElement, filter.SelectFields)
			if err != nil {
				return err
			}

			if err = sw.SetRow(axis, cells); err != nil {
				return errors.NoType.Wrap(err, "error setting sheet row")
			}
		}
	}

	return nil
}

// getCells получение ячеек для записи в xlsx
func getCells(model *focus.Model, modelElement any, selectedFields []string) ([]any, error) {
	fieldsMap, err := model.ElementToMap(modelElement, func(field focus.Field) bool {
		return slices.Contains(selectedFields, field.Code)
	})
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error convert model element to map")
	}

	var cells = make([]any, len(model.Fields))
	for i, fieldCode := range selectedFields {
		cells[i] = excelize.Cell{Value: fieldsMap[fieldCode]}
	}

	return cells, nil
}
