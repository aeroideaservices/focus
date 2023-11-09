package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/aeroideaservices/focus/configurations/plugin/actions"
	"github.com/aeroideaservices/focus/configurations/plugin/entity"
	"github.com/aeroideaservices/focus/configurations/postgres/service/util"
	"github.com/google/uuid"
	stackedErrors "github.com/pkg/errors"
	"gorm.io/gorm"
)

type confRepository struct {
	db *gorm.DB
}

// NewConfigurationRepository конструктор
func NewConfigurationRepository(db *gorm.DB) actions.ConfigurationsRepository {
	return &confRepository{db: db}
}

// Get получение конфигурации по id
func (r confRepository) Get(ctx context.Context, Id uuid.UUID) (*entity.Configuration, error) {
	confEntity := &entity.Configuration{}
	err := r.db.WithContext(ctx).
		First(confEntity, Id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, actions.ErrConfNotFound
	}
	if err != nil {
		return nil, stackedErrors.WithStack(err)
	}

	return confEntity, nil
}

// Has проверка существования конфигурации по id
func (r confRepository) Has(ctx context.Context, id uuid.UUID) bool {
	err := r.db.WithContext(ctx).
		Select("id").
		First(&entity.Configuration{}, "id = ?", id).
		Error

	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// HasByCode проверка существования конфигурации по коду
func (r confRepository) HasByCode(ctx context.Context, code string) bool {
	db := r.db
	confEntity := &entity.Configuration{}
	db.WithContext(ctx).First(confEntity, "code = ?", code)
	return confEntity.Id != uuid.Nil
}

// Create создание конфигурации
func (r confRepository) Create(ctx context.Context, conf ...entity.Configuration) error {
	err := r.db.WithContext(ctx).Create(conf).Error
	return stackedErrors.WithStack(err)
}

// Update обновление конфигурации
func (r confRepository) Update(ctx context.Context, conf *entity.Configuration) error {
	err := r.db.WithContext(ctx).Save(conf).Error
	return stackedErrors.WithStack(err)
}

// FindByCode получение конфигурации по коду
func (r confRepository) FindByCode(ctx context.Context, code string) (*entity.Configuration, error) {
	confEntity := &entity.Configuration{}
	err := r.db.WithContext(ctx).First(confEntity, "code = ?", code).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, actions.ErrConfNotFound
	}
	if err != nil {
		return nil, stackedErrors.WithStack(err)
	}

	return confEntity, nil
}

// Delete удаление конфигурации
func (r confRepository) Delete(ctx context.Context, Id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&entity.Configuration{Id: Id}).Error
	return stackedErrors.WithStack(err)
}

// List получение списка конфигураций по фильтру
func (r confRepository) List(ctx context.Context, dto actions.ConfigurationsListFilter) ([]entity.Configuration, error) {
	var confEntityList = make([]entity.Configuration, 0)

	sort := dto.Sort
	if sort == "" {
		sort = "updated_at"
	}
	order := dto.Order
	if order == "" {
		order = "desc"
	}

	db := r.db.WithContext(ctx)
	db = filterConf(db, dto.Filter)
	err := db.
		Limit(dto.Limit).Offset(dto.Offset).
		Order(fmt.Sprintf("%s %s", sort, order)).
		Find(&confEntityList).Error
	if err != nil {
		return nil, stackedErrors.WithStack(err)
	}

	return confEntityList, nil
}

// Count получение количества конфигураций по фильтру
func (r confRepository) Count(ctx context.Context, filter actions.ConfigurationsFilter) (int, error) {
	var count int64
	db := r.db.WithContext(ctx).Model(&entity.Configuration{})
	db = filterConf(db, filter)
	err := db.Count(&count).Error
	if err != nil {
		return 0, stackedErrors.WithStack(err)
	}

	return int(count), nil
}

func filterConf(db *gorm.DB, filter actions.ConfigurationsFilter) *gorm.DB {
	if value := filter.Query; value != "" {
		util.PrepareLikeValue(&value)
		db = db.Where("id::varchar ILIKE ? OR code ILIKE ? or name ILIKE ?", value, value, value)
	}

	return db
}
