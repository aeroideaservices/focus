package entity

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	Id        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	Name      string     `json:"name"`
	Filename  string     `json:"filename"`
	Alt       string     `json:"alt"`
	Title     string     `json:"title"`
	Size      int64      `json:"size"`
	Filepath  string     `json:"filepath" gorm:"unique_index"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	FolderId  *uuid.UUID `json:"folderId" gorm:"type:uuid"`
}

type MediaList struct {
	Total int     `json:"total"`
	Items []Media `json:"items"`
}

func (Media) TableName() string {
	return "media"
}
