package entity

import "github.com/google/uuid"

//type ExternalTag struct {
//	ID   uuid.UUID `focus:"primaryKey;code:id;column:id;title:ID" validate:"required,notBlank"`
//	Text string    `focus:"title:Текст;filterable" validate:"required,min=1,max=50"`
//	Link string    `focus:"title:Ссылка;filterable" validate:"required,min=1,max=50"`
//}
//
//func (ExternalTag) TableName() string {
//	return "external_tags"
//}
//
//type InternalTag struct {
//	ID          uuid.UUID `focus:"primaryKey;code:id;column:id;title:ID" validate:"required,notBlank"`
//	Text        string    `focus:"title:Текст;filterable" validate:"required,min=1,max=50"`
//	PageCode    string    `focus:"title:Код страницы;filterable" validate:"required,min=1,max=50"`
//	GalleryCode string    `focus:"title:Код галереи;filterable" validate:"required,min=1,max=50"`
//}

type Tag struct {
	ID   uuid.UUID `focus:"primaryKey;code:id;column:id;title:ID" validate:"required,notBlank"`
	Text string    `focus:"title:Текст;filterable" validate:"required,min=1,max=50"`
	Link string    `focus:"title:Ссылка;filterable" validate:"required"`
}

func (Tag) TableName() string {
	return "tags"
}

func (Tag) ModelTitle() string {
	return "Теги"
}

//func (InternalTag) TableName() string {
//	return "internal_tags"
//}
