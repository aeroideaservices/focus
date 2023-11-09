package entity

import (
	"github.com/google/uuid"
	"time"
)

// Configuration сущность конфигурации
type Configuration struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type ConfigurationsList struct {
	Total int             `json:"total"`
	Items []Configuration `json:"items"`
}

func (Configuration) TableName() string {
	return "configurations"
}
