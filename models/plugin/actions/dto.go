package actions

import (
	"github.com/aeroideaservices/focus/models/plugin/form"
	"io"
)

type CreateModelElement struct {
	ModelCode    string         `validate:"required"`
	ModelElement map[string]any `validate:"required"`
}

type UpdateModelElement struct {
	ModelCode    string         `validate:"required"`
	PKey         any            `validate:"required"`
	ModelElement map[string]any `validate:"required"`
}

type DeleteModelElement struct {
	ModelCode string `json:"modelCode" validate:"required"`
	PKey      any    `json:"pKey" validate:"required,notBlank"`
}

type DeleteModelElements struct {
	ModelCode string `json:"modelCode" validate:"required"`
	PKeys     []any  `json:"pKeys" validate:"required,unique"`
}

type GetModelElement struct {
	ModelCode string `json:"modelCode" validate:"required,notBlank"`
	PKey      any    `json:"PKey" validate:"required,notBlank"`
}

type ListModelElementsQuery struct {
	ModelCode    string              `json:"-"`
	Filter       ModelElementsFilter `json:"filter"`
	SelectFields []string            `json:"-"`
	Pagination   `json:"-"`
	OrderBy      `json:"-"`
}

type ListModelElements struct {
	ModelCode string       `json:"-"`
	Filter    FieldsFilter `json:"filter"`
	ModelElementsQueryFilter
	Pagination
	OrderBy
}

type Select []string

type List struct {
	Items []map[string]any `json:"items"`
	Total int64            `json:"total"`
}

type FieldsFilter map[string][]any

type ModelElementsFilter struct {
	FieldsFilter FieldsFilter             `json:"fieldsFilter"`
	QueryFilter  ModelElementsQueryFilter `json:"query"`
}

type ModelElementsQueryFilter struct {
	FieldsCodes []string `json:"fields" validate:"required_with=Query"`
	Query       string   `json:"query"`
}

type OrderBy struct {
	Sort  string `json:"sort"`
	Order string `json:"order" validate:"omitempty,oneof=asc desc"`
}

type Pagination struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"required,min=1,max=100"`
}

type ModelDescription struct {
	Code           string     `json:"code"`           // Уникальный код модели
	Title          string     `json:"name"`           // Название модели
	IdentifierCode string     `json:"identifierCode"` // Код поля, которое является идентификатором модели
	Views          ModelViews `json:"views"`          // Параметры отображения полей модели
}

type ModelsList struct {
	Items []ModelShort `json:"items"`
	Total int          `json:"total"`
}

type ModelShort struct {
	Code  string `json:"code"` // Уникальный код модели
	Title string `json:"name"` // Название модели
}

// ModelViews описывает параметры отображения полей модели
type ModelViews struct {
	Create EditView `json:"create"`
	Update EditView `json:"update"`
	Filter FormView `json:"filter"`
	List   ListView `json:"list"`
}

type EditView struct {
	FormFields []FormField `json:"formFields"`
	Validation any         `json:"validation"`
}

type FormView struct {
	FormFields []FormField `json:"formFields"`
}

type ListView struct {
	Fields []ListField `json:"fields"`
}

type ListField struct {
	Code     string `json:"code"`
	Title    string `json:"name"`
	Sortable bool   `json:"sortable"`
	IsTime   bool   `json:"isTime"`
}

type FieldValuesList struct {
	Items any   `json:"items"`
	Total int64 `json:"total"`
}

type FormField struct {
	Code      string         `json:"code"`
	Title     string         `json:"name"`
	Type      form.FieldType `json:"type"`
	Multiple  bool           `json:"multiple"`
	Sortable  bool           `json:"sortable"`
	Block     string         `json:"block,omitempty"`
	Extra     map[string]any `json:"extra,omitempty"`
	Hidden    bool           `json:"hidden,omitempty"`
	Disabled  bool           `json:"disabled,omitempty"`
	Step      float64        `json:"step,omitempty"`
	Precision int            `json:"precision,omitempty"`
}

type GetModel struct {
	ModelCode string `json:"modelCode" validate:"required"`
}

type ListModels struct {
	Pagination
	ModelOrderBy
}

type ListFieldValues struct {
	ModelCode string `json:"modelCode" validate:"required"`
	FieldCode string `json:"fieldCode" validate:"required"`
	Query     string `json:"query"`
	Pagination
}

type ModelOrderBy struct {
	Sort  string `json:"sort" validate:"omitempty,oneof=code name"`
	Order string `json:"order" validate:"omitempty,oneof=asc desc"`
}

type ExportModelElements struct {
	ModelCode string       `json:"modelCode" validate:"required"`
	Filter    FieldsFilter `json:"filter"`
	OrderBy
}

type GetExportInfo struct {
	ModelCode string
}

type CreateFile struct {
	Key         string
	ContentType string
	File        io.ReadCloser
}
