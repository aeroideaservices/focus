package actions

import (
	"context"
	"github.com/aeroideaservices/focus/configurations/plugin/entity"
	"github.com/google/uuid"
)

// ConfigurationsRepository интерфейс репозитория конфигураций
type ConfigurationsRepository interface {
	Get(ctx context.Context, id uuid.UUID) (conf *entity.Configuration, err error)
	Has(ctx context.Context, id uuid.UUID) bool
	HasByCode(ctx context.Context, code string) bool
	Create(ctx context.Context, conf ...entity.Configuration) error
	Update(ctx context.Context, conf *entity.Configuration) error
	FindByCode(ctx context.Context, code string) (conf *entity.Configuration, err error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, dto ConfigurationsListFilter) ([]entity.Configuration, error)
	Count(ctx context.Context, filter ConfigurationsFilter) (int, error)
}

type ConfigurationsListFilter struct {
	Offset int
	Limit  int
	Sort   string
	Order  string
	Filter ConfigurationsFilter
}

type ConfigurationsFilter struct {
	Query string
}

// OptionsRepository интерфейс репозитория настроек
type OptionsRepository interface {
	Has(ctx context.Context, confId uuid.UUID, optId uuid.UUID) bool
	HasByCode(ctx context.Context, confId uuid.UUID, optCodes ...string) bool
	Get(ctx context.Context, id uuid.UUID) (*entity.Option, error)
	Create(ctx context.Context, opt ...entity.Option) error
	Update(ctx context.Context, opt ...entity.Option) error
	Delete(ctx context.Context, optId uuid.UUID) error
	List(ctx context.Context, filter OptionsListFilter) ([]entity.Option, error)
	ListPreviews(ctx context.Context, filter OptionsListFilter) (optShorts []entity.OptionShort, err error)
	Count(ctx context.Context, filter OptionsFilter) (int, error)
}

type OptionsListFilter struct {
	Offset int
	Limit  int
	Sort   string
	Order  string
	Filter OptionsFilter
}

type OptionsFilter struct {
	ConfId   uuid.UUID
	ConfCode string
	OptCodes []string
}

// MediaService интерфейс сервиса медиа
type MediaService interface {
	CheckIds(ctx context.Context, ids ...uuid.UUID) error
}
