package entity

import (
	"github.com/google/uuid"
	"time"
)

type exportStatus string

var (
	StatusPending exportStatus = "pending"
	StatusSucceed exportStatus = "succeed"
	StatusError   exportStatus = "error"
)

type ExportInfo struct {
	ID        uuid.UUID    `json:"id" gorm:"primaryKey;type:uuid"`
	ModelCode string       `json:"modelCode"`
	Filepath  string       `json:"filepath"`
	Status    exportStatus `json:"status"`
	Time      time.Time    `json:"time"`
}

func (ExportInfo) TableName() string {
	return "model_export_info"
}
