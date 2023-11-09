package actions

import (
	"github.com/aeroideaservices/focus/configurations/plugin/entity"
	"github.com/google/uuid"
)

type CreateConfiguration struct {
	Code string `json:"code" validate:"required,min=3,max=50,sluggable"`
	Name string `json:"name" validate:"required,notBlank,min=3,max=50"`
}

type GetConfiguration struct {
	Id uuid.UUID `json:"id" validate:"required,notBlank"`
}

type UpdateConfiguration struct {
	Id   uuid.UUID `json:"id" validate:"required,notBlank"`
	Code string    `json:"code" validate:"required"`
	Name string    `json:"name" validate:"required,notBlank,min=3,max=50"`
}

type ListConfigurations struct {
	Offset int    `json:"offset" validate:"min=0"`
	Limit  int    `json:"limit" validate:"required,min=10,max=100"`
	Sort   string `json:"sort" validate:"omitempty,oneof=id code name"`
	Order  string `json:"order" validate:"omitempty,oneof=asc desc"`
	Query  string `json:"query" validate:""`
}

type CreateOption struct {
	ConfId uuid.UUID `json:"confId" validate:"required,notBlank"`
	Code   string    `json:"code" validate:"required,min=3,max=50,sluggable"`
	Name   string    `json:"name" validate:"required,notBlank,min=3,max=50"`
	Type   string    `json:"type" validate:"required,oneof=string text integer checkbox file image datetime date"`
}

type UpdateOption struct {
	Id     uuid.UUID `json:"id" validate:"required,notBlank"`
	ConfId uuid.UUID `json:"confId" validate:"required,notBlank"`
	Code   string    `json:"code" validate:"required"`
	Name   string    `json:"name" validate:"required,notBlank,min=3,max=50"`
	Type   string    `json:"type" validate:"required,oneof=string text integer checkbox file image datetime date"`
}

type UpdateOptionsList struct {
	Items  []entity.OptionShort `json:"items" validate:"required,min=1,unique=Code"`
	ConfId uuid.UUID            `json:"confId" validate:"required,notBlank"`
}

type GetOption struct {
	Id     uuid.UUID `json:"id" validate:"required,notBlank"`
	ConfId uuid.UUID `json:"confId" validate:"required,notBlank"`
}

type ListOptionsPreviews struct {
	ConfCode string   `json:"confCode" validate:"required,min=3,max=50"`
	OptCodes []string `json:"optCodes" validate:"omitempty,unique,gt=0,dive,min=3,max=50"`
}

type ListOptions struct {
	ConfId uuid.UUID `json:"confId" validate:"required,notBlank"`
	Offset int       `json:"offset" validate:"min=0"`
	Limit  int       `json:"limit" validate:"required,min=10,max=100"`
	Sort   string    `json:"sort"`
	Order  string    `json:"order" validate:"omitempty,oneof=asc desc"`
}
