package entity

import (
	"github.com/aeroideaservices/focus/media/plugin/entity"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID     `focus:"primaryKey;code:id;column:id;title:ID" validate:"required,notBlank"`
	FirstName string        `focus:"title:Имя;filterable" validate:"required,min=1,max=50"`
	LastName  string        `focus:"title:Фамилия;filterable" validate:"required,min=1,max=50"`
	Position  string        `focus:"title:Должность;filterable" validate:"required,min=1,max=50"`
	PictureId *uuid.UUID    `focus:"-" validate:"-"`
	Picture   *entity.Media `focus:"title:Картинки;media;hidden:list"`
}

func (User) TableName() string {
	return "users"
}
