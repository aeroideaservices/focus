package entity

import "github.com/google/uuid"

// Menu сущность меню
type Menu struct {
	Id   uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

func (Menu) TableName() string {
	return "menus"
}
