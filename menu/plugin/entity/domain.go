package entity

import "github.com/google/uuid"

type Domain struct {
	Id     uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Domain string    `json:"domain"`
}

func (Domain) TableName() string {
	return "domains"
}
