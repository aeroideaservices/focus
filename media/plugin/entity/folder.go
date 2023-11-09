package entity

import (
	"time"

	"github.com/google/uuid"
)

type Folder struct {
	Id        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	FolderId  *uuid.UUID `json:"parentFolderId"`
	Folder    *Folder    `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Medias    []Media    `gorm:"constraint:OnDelete:CASCADE" json:"-"`
}

func (Folder) TableName() string {
	return "folders"
}
