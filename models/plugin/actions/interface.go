package actions

import (
	"context"
	"github.com/aeroideaservices/focus/models/plugin/entity"
	"github.com/aeroideaservices/focus/models/plugin/focus"
	"github.com/google/uuid"
	"os"
)

type ModelsRegistry interface {
	GetModel(code string) *focus.Model
	ListModels() []*focus.Model
}

type Repository interface {
	Has(ctx context.Context, pk any) (bool, error)
	Create(ctx context.Context, elem any) (id any, err error)
	Get(ctx context.Context, key any) (elem any, err error)
	Update(ctx context.Context, elem any) error
	Count(ctx context.Context, filter ModelElementsFilter) (count int64, err error)
	List(ctx context.Context, filter ListModelElementsQuery) (elems []any, err error)
	ListFieldValues(ctx context.Context, action ListFieldValues) (fieldValues any, err error)
	CountFieldValues(ctx context.Context, code string, query string) (count int64, err error)
	Delete(ctx context.Context, pks ...any) error
}

type RepositoryResolver interface {
	Resolve(modelCode string) Repository
}

type MediaService interface {
	CheckIds(ctx context.Context, ids ...uuid.UUID) error
}

type Validator interface {
	Validate(ctx context.Context, value any) error
	ValidatePartial(ctx context.Context, value any) error
}

type Exporter interface {
	GetFile(ctx context.Context, model *focus.Model, filter ListModelElementsQuery) (*os.File, error)
}

type FileStorage interface {
	Upload(ctx context.Context, media *CreateFile) error
	Delete(ctx context.Context, keys ...string) error
	GetSize(ctx context.Context, key string) (int64, error)
}

type ExportInfoRepository interface {
	GetLast(ctx context.Context, code string) (*entity.ExportInfo, error)
	Create(ctx context.Context, export entity.ExportInfo) error
	Update(ctx context.Context, export entity.ExportInfo) error
	Delete(ctx context.Context, id uuid.UUID) (string, error)
}
