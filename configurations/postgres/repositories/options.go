package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/aeroideaservices/focus/configurations/plugin/actions"
	"github.com/aeroideaservices/focus/configurations/plugin/entity"
	"github.com/google/uuid"
	stackedErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type optionRepository struct {
	db *gorm.DB
}

// NewOptionsRepository конструктор
func NewOptionsRepository(db *gorm.DB) actions.OptionsRepository {
	return &optionRepository{db: db}
}

// Has проверка существования настройки по id
func (r optionRepository) Has(ctx context.Context, confId uuid.UUID, optId uuid.UUID) bool {
	err := r.db.WithContext(ctx).
		Select("options.id").
		Where("options.conf_id = ?", confId).
		Where("options.id = ?", optId).
		First(&entity.Option{}).Error

	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// HasByCode проверка существования настройки по коду
func (r optionRepository) HasByCode(ctx context.Context, confId uuid.UUID, optCodes ...string) bool {
	err := r.db.WithContext(ctx).
		Select("options.id").
		Where("options.conf_id = ?", confId).
		Where("options.code IN (?)", optCodes).
		First(&entity.Option{}).Error

	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Get получение настройки по id
func (r optionRepository) Get(ctx context.Context, id uuid.UUID) (*entity.Option, error) {
	opt := &entity.Option{}
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(opt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, actions.ErrOptNotFound
	}
	if err != nil {
		return nil, stackedErrors.WithStack(err)
	}

	return opt, nil
}

// Create создание настройки
func (r optionRepository) Create(ctx context.Context, opt ...entity.Option) error {
	err := r.db.WithContext(ctx).Create(opt).Error
	if err != nil {
		return stackedErrors.WithStack(err)
	}

	return nil
}

// Update обновление настройки
func (r optionRepository) Update(ctx context.Context, opt ...entity.Option) error {
	err := r.db.WithContext(ctx).Save(&opt).Error

	return stackedErrors.WithStack(err)
}

// Delete удаление настройки
func (r optionRepository) Delete(ctx context.Context, optId uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&entity.Option{Id: optId}).Error
	if err != nil {
		return stackedErrors.WithStack(err)
	}

	return nil
}

// List получение списка настроек по фильтру
func (r optionRepository) List(ctx context.Context, filter actions.OptionsListFilter) ([]entity.Option, error) {
	opts := make([]entity.Option, 0)
	err := r.getList(ctx, filter, &opts)
	if err != nil {
		return nil, stackedErrors.WithStack(err)
	}

	return opts, nil
}

// ListPreviews получение списка превью настроек по фильтру
func (r optionRepository) ListPreviews(ctx context.Context, filter actions.OptionsListFilter) ([]entity.OptionShort, error) {
	options := make([]entity.OptionShort, 0)
	err := r.getList(ctx, filter, &options)
	if err != nil {
		return nil, stackedErrors.WithStack(err)
	}

	return options, nil
}

func (r optionRepository) getList(ctx context.Context, filter actions.OptionsListFilter, output interface{}) error {
	sort := "options." + filter.Sort
	if filter.Sort == "" {
		sort = "options.updated_at"
	}
	order := filter.Order
	if order == "" {
		order = "desc"
	}

	db := r.db.WithContext(ctx).Model(entity.Option{}).Distinct()
	db = r.filterOpt(db, filter.Filter)
	if filter.Limit != 0 {
		db = db.Limit(filter.Limit)
	}
	err := db.Offset(filter.Offset).
		Order(fmt.Sprintf("%s %s", sort, order)).
		Scan(output).Error

	return stackedErrors.WithStack(err)
}

// Count получение количества настроек по фильтру
func (r optionRepository) Count(ctx context.Context, filter actions.OptionsFilter) (int, error) {
	var count int64

	db := r.db.WithContext(ctx).Model(entity.Option{}).Distinct("options.id")
	db = r.filterOpt(db, filter)
	err := db.Count(&count).Error
	if err != nil {
		return 0, stackedErrors.WithStack(err)
	}

	return int(count), nil
}

func (optionRepository) filterOpt(db *gorm.DB, filter actions.OptionsFilter) *gorm.DB {
	if filter.ConfId != uuid.Nil {
		db = db.Where("options.conf_id = ?", filter.ConfId)
	}
	if filter.ConfCode != "" {
		db = db.Joins("INNER JOIN configurations conf on conf.id = options.conf_id AND conf.code = ?", filter.ConfCode)
	}
	if len(filter.OptCodes) > 0 {
		db = db.Where("options.code IN(?)", filter.OptCodes)
	}

	return db
}
