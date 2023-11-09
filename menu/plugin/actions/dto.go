package actions

import (
	"github.com/aeroideaservices/focus/menu/plugin/entity"
	"github.com/google/uuid"
)

type ListMenus struct {
	Pagination
	Sort  string `json:"sort" validate:"omitempty,oneof=name"`
	Order string `json:"order" validate:"omitempty,oneof=asc desc"`
}

type MenusList struct {
	Total int64         `json:"total"`
	Items []entity.Menu `json:"items"`
}

type GetMenu struct {
	MenuId uuid.UUID `json:"menuId" validate:"required,notBlank"`
}

type GetMenuByCode struct {
	MenuCode string `json:"menuCode" validate:"required,min=3,max=50"`
}

type CreateMenu struct {
	Name string `json:"name" validate:"required,notBlank,min=3,max=50"`
	Code string `json:"code" validate:"required,min=3,max=50,sluggable"`
}

type UpdateMenu struct {
	Id   uuid.UUID `json:"id" validate:"required,notBlank"`
	Name string    `json:"name" validate:"required,notBlank,min=3,max=50"`
	Code string    `json:"code" validate:"required,min=3,max=50,sluggable"`
}

type ListMenuItems struct {
	Sort   string         `json:"sort" validate:"omitempty,oneof=id name position"`
	Order  string         `json:"order" validate:"omitempty,oneof=asc desc"`
	Filter MenuItemFilter `json:"filter"`
}

type CreateMenuItem struct {
	MenuId           uuid.UUID         `json:"menuId" validate:"required,notBlank"`
	Name             string            `json:"name" validate:"required,notBlank,min=3,max=50"`
	DomainId         *uuid.UUID        `json:"domainId" validate:"omitempty,notBlank"`
	Url              string            `json:"url" validate:"omitempty,slashedSluggable"`
	AdditionalFields []AdditionalField `json:"additionalFields" validate:"omitempty,unique=Code,dive"`
	ParentMenuItemId *uuid.UUID        `json:"parentMenuItemId" validate:"omitempty,notBlank"`
}

type AdditionalField struct {
	Code  string `json:"code" validate:"required,min=3,max=50"`
	Value string `json:"value" validate:"required,min=1"`
}

type GetMenuItem struct {
	MenuId     uuid.UUID `json:"menuId" validate:"required,notBlank"`
	MenuItemId uuid.UUID `json:"menuItemId" validate:"required,notBlank"`
}

type UpdateMenuItem struct {
	MenuId           uuid.UUID         `json:"menuId" validate:"required,notBlank"`
	MenuItemId       uuid.UUID         `json:"menuItemId" validate:"required,notBlank"`
	Name             string            `json:"name" validate:"required,notBlank,min=3,max=50"`
	DomainId         *uuid.UUID        `json:"domainId" validate:"omitempty,notBlank"`
	Url              string            `json:"url" validate:"omitempty,slashedSluggable"`
	AdditionalFields []AdditionalField `json:"additionalFields" validate:"omitempty,unique=Code"`
}

type MoveMenuItem struct {
	MenuId           uuid.UUID  `json:"menuId" validate:"required,notBlank"`
	MenuItemId       uuid.UUID  `json:"menuItemId" validate:"required,notBlank"`
	ParentMenuItemId *uuid.UUID `json:"parentMenuItemId" validate:"omitempty,notBlank,nefield=MenuItemId"`
	Position         int64      `json:"position" validate:"required,min=1"`
}

type Pagination struct {
	Offset int `json:"offset" validate:"min=0"`
	Limit  int `json:"limit" validate:"required,min=10,max=100"`
}

type ListDomains struct {
	Pagination
}

type DomainsList struct {
	Total int64           `json:"total"`
	Items []entity.Domain `json:"items"`
}

type CreateDomain struct {
	Domain string `json:"domain" validate:"required,url"`
}
