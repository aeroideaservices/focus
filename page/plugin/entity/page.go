package entity

import (
	"github.com/google/uuid"
)

type Page struct {
	ID             uuid.UUID        `gorm:"type:uuid;primaryKey;" json:"id"`
	Name           string           `json:"name"`
	Code           string           `json:"code"`
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	TitleSeo       string           `json:"titleSeo"`
	DescriptionSeo string           `json:"descriptionSeo"`
	Keywords       string           `json:"keywords"`
	IsPublished    bool             `json:"isPublished"`
	Sort           int              `json:"sort"`
	PagesGalleries []PagesGalleries `json:"pagesGalleries" gorm:"foreignKey:PagesID"`

	OgType string `json:"ogType"`
}

func (Page) TableName() string {
	return "pages"
}

type PagesGalleries struct {
	PagesID   *uuid.UUID
	GalleryID *uuid.UUID
	Gallery   Gallery `gorm:"foreignKey:GalleryID"`
	Position  int
}
