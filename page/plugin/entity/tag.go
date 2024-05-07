package entity

import "github.com/google/uuid"

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
