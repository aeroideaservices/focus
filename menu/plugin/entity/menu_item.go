package entity

import (
	"github.com/aeroideaservices/focus/services/db/db_types/json"
	"github.com/google/uuid"
)

// MenuItem сущность пункта меню
type MenuItem struct {
	Id               uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Name             string     `json:"name"`
	Domain           *Domain    `json:"-"`
	DomainId         *uuid.UUID `json:"domainId"`
	Url              string     `json:"url"`
	Position         int64      `json:"position"`
	AdditionalFields json.JSONB `json:"additionalFields"`
	ParentMenuItemId *uuid.UUID `json:"parentMenuItemId" gorm:"default:NULL"`
	ParentMenuItem   *MenuItem  `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	MenuId           uuid.UUID  `json:"menuId"`
	Menu             Menu       `json:"-" gorm:"constraint:OnDelete:CASCADE"`
}

func (MenuItem) TableName() string {
	return "menu_items"
}
