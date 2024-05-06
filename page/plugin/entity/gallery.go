package entity

import (
	"github.com/google/uuid"
)

type Gallery struct {
	ID             uuid.UUID        `focus:"primaryKey;code:id;column:id;title:ID" validate:"required,notBlank"`
	Name           string           `focus:"title:Название;filterable" validate:"required,min=1,max=50"`
	Code           string           `focus:"title:Код;unique;sluggableOn:name;disabled:update;sortable;filterable;viewExtra:sluggableOnName" validate:"required,notBlank,sluggable"`
	GalleriesCards []GalleriesCards `focus:"title:Карточки;many2many:galleries_cards;viewExtra:selectCards" gorm:"foreignKey:GalleryID" validate:"omitempty,unique=ID,dive,notBlank,structonly"`
	Hidden         bool             `focus:"title:Скрыт в меню;" validate:"required"`
	IsPublished    bool             `focus:"title:Опубликована ли галерея;" validate:"required"`
}

func (Gallery) TableName() string {
	return "galleries"
}

func (Gallery) ModelTitle() string {
	return "Галереи"
}

type GalleriesCards struct {
	GalleryID *uuid.UUID `focus:"primaryKey;code:galleryId;column:gallery_id;title:ID" validate:"required,notBlank"`
	CardID    *uuid.UUID `focus:"primaryKey;code:cardId;column:card_id;title:ID" validate:"required,notBlank"`
	Card      *Card      `focus:"title:Карточки;view:select;hidden:list" gorm:"foreignKey:CardID"`
	Position  int        `focus:"-" validate:"-"`
}

func (GalleriesCards) TableName() string {
	return "galleries_cards"
}

func (GalleriesCards) ModelTitle() string {
	return "Галереи"
}
