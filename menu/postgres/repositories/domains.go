package repositories

import (
	"context"
	"github.com/aeroideaservices/focus/menu/plugin/actions"
	"github.com/aeroideaservices/focus/menu/plugin/entity"
	"github.com/aeroideaservices/focus/services/errors"
	"gorm.io/gorm"
)

type domainRepository struct {
	db *gorm.DB
}

// NewDomainRepository конструктор
func NewDomainRepository(db *gorm.DB) actions.DomainsRepository {
	return &domainRepository{db: db}
}

// Has проверка существования домена
func (r domainRepository) Has(ctx context.Context, domain string) (bool, error) {
	err := r.db.WithContext(ctx).
		Where("domain", domain).
		First(&entity.Domain{}).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, errors.NoType.Wrap(err, "error checking if domain exists")
	}

	return true, nil
}

// Create создание домена
func (r domainRepository) Create(ctx context.Context, domain entity.Domain) error {
	err := r.db.WithContext(ctx).
		Create(&domain).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error creating domain")
	}

	return nil
}

// List получение списка доменов
func (r domainRepository) List(ctx context.Context, query actions.DomainsListQuery) ([]entity.Domain, error) {
	var domains []entity.Domain

	err := r.db.WithContext(ctx).
		Limit(query.Limit).Offset(query.Offset).
		Find(&domains).Error
	if err != nil {
		return nil, errors.Internal.Wrap(err, "error listing domain")
	}

	return domains, nil
}

// Count получение общего количества доменов
func (r domainRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&entity.Domain{}).
		Count(&count).Error
	if err != nil {
		return 0, errors.Internal.Wrap(err, "error counting domain")
	}

	return count, nil
}
