package postgres

import (
	"context"
	"github.com/aeroideaservices/focus/models/plugin/entity"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type exportInfoRepository struct {
	db *gorm.DB
}

// NewExportInfoRepository конструктор
func NewExportInfoRepository(db *gorm.DB) *exportInfoRepository {
	return &exportInfoRepository{db: db}
}

// GetLast получение инфо по последнему экспорту
func (r exportInfoRepository) GetLast(ctx context.Context, modelCode string) (*entity.ExportInfo, error) {
	exportInfo := &entity.ExportInfo{}
	err := r.db.WithContext(ctx).Where("model_code", modelCode).Order("time desc").First(exportInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrap(err, "export info not found")
	}
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting last export info")
	}

	return exportInfo, nil
}

// Create создание новой записи об экспорте
func (r exportInfoRepository) Create(ctx context.Context, info entity.ExportInfo) error {
	err := r.db.WithContext(ctx).Create(info).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error creating export info")
	}

	return nil
}

// Update обновление инфо по экспорту
func (r exportInfoRepository) Update(ctx context.Context, info entity.ExportInfo) error {
	err := r.db.WithContext(ctx).Save(info).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error updating export info")
	}

	return nil
}

// Delete удаление инфо по экспорту
func (r exportInfoRepository) Delete(ctx context.Context, id uuid.UUID) (string, error) {
	exportInfo := &entity.ExportInfo{}
	err := r.db.WithContext(ctx).Where(id).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "filepath"}}}).
		Delete(exportInfo).Error
	if err != nil {
		return "", errors.NoType.Wrap(err, "error deleting export info")
	}

	return exportInfo.Filepath, nil
}
