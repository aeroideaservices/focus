package entity

import (
	"github.com/aeroideaservices/focus/media/plugin/entity"
	"github.com/google/uuid"
	entity2 "gitlab.aeroidea.ru/internal-projects/focus/forms/plugin/entity"
)

type Card struct {
	ID            uuid.UUID    `focus:"primaryKey;code:id;column:id;title:ID" validate:"required,notBlank"`
	Name          string       `focus:"title:Название"`
	Code          string       `focus:"title:Код"`
	Type          string       `focus:"title:Название;filterable" validate:"required,min=1,max=50"`
	IsPublished   bool         `focus:"title:Опубликована ли карточка;" validate:"required"`
	RegularCard   *RegularCard `focus:"title:Привязанная обычная карточка;view:select;viewExtra:regularCardSelect;hidden:list" validate:"structonly"`
	RegularCardId *uuid.UUID   `focus:"-" validate:"-"`
	VideoCard     *VideoCard   `focus:"title:Привязанная карточка с видео;view:select;viewExtra:videoCardSelect;hidden:list" validate:"structonly"`
	VideoCardId   *uuid.UUID   `focus:"-" validate:"-"`
	HtmlCard      *HtmlCard    `focus:"title:Привязанные html контейнеры;view:select;viewExtra:htmlCardSelect;hidden:list" validate:"structonly"`
	HtmlCardId    *uuid.UUID   `focus:"-" validate:"-"`
	PhotoCard     *PhotoCard   `focus:"title:Привязанная карточка с фото;view:select;viewExtra:photoCardSelect;hidden:list" validate:"structonly"`
	PhotoCardId   *uuid.UUID   `focus:"-" validate:"-"`
	FormCard      *FormCard    `focus:"title:Привязанная карточка с фото;view:select;hidden:list" validate:"structonly"`
	FormCardId    *uuid.UUID   `focus:"-" validate:"-"`

	Title       string `focus:"title:Заголовок"`
	Description string `focus:"title:Описание"`
	OgType      string `focus:"title:Тип для og;filterable"`
}

func (Card) TableName() string {
	return "cards"
}

type RegularCard struct {
	ID                 uuid.UUID          `focus:"primaryKey;code:id;column:id;title:ID" validate:"required,notBlank"`
	PreviewText        string             `focus:"title:Превью текст" validate:"required"`
	DetailText         string             `focus:"title:Весь текст" validate:"required"`
	Inverted           bool               `json:"inverted"`
	Video              *entity.Media      `focus:"title:Видео;media;hidden:list"`
	VideoId            *uuid.UUID         `focus:"-"`
	VideoLite          *entity.Media      `focus:"title:Легкое видео;media;hidden:list"`
	VideoLiteId        *uuid.UUID         `focus:"-"`
	VideoPreview       *entity.Media      `focus:"title:Видео превью;media;hidden:list"`
	VideoPreviewId     *uuid.UUID         `focus:"-"`
	VideoPreviewBlur   *entity.Media      `focus:"title:Видео превью с блюром;media;hidden:list"`
	VideoPreviewBlurId *uuid.UUID         `focus:"-"`
	RegularCardsTags   []RegularCardsTags `focus:"-" gorm:"foreignKey:RegularCardID" validate:"omitempty,unique=ID,dive,notBlank,structonly"`
	LearnMoreUrl       *string            `focus:"title:Ссылка узнать больше"`
	User               *User              `focus:"title:Спикер;view:select;viewExtra:userSelect;hidden:list" validate:"omitempty,structonly"`
	UserId             *uuid.UUID         `focus:"-"`
}

func (RegularCard) TableName() string {
	return "regular_cards"
}

type RegularCardsTags struct {
	RegularCardID uuid.UUID `focus:"primaryKey;code:regularCardId;column:regular_card_id;title:ID" validate:"required,notBlank"`
	TagID         uuid.UUID `focus:"primaryKey;code:intTagId;column:int_tag_id;title:ID" validate:"required,notBlank"`
	Tag           Tag       `focus:"title:Тег;view:select;hidden:list" gorm:"foreignKey:TagID"`
}

func (RegularCardsTags) TableName() string {
	return "regular_cards_tags"
}

type VideoCard struct {
	ID                 uuid.UUID     `focus:"primaryKey;code:id;column:id;title:ID" validate:"required,notBlank"`
	Video              *entity.Media `focus:"title:Видео;media;hidden:list"`
	VideoId            *uuid.UUID    `focus:"-"`
	VideoLite          *entity.Media `focus:"title:Легкое видео;media;hidden:list"`
	VideoLiteId        *uuid.UUID    `focus:"-"`
	VideoPreview       *entity.Media `focus:"title:Видео превью;media;hidden:list"`
	VideoPreviewId     *uuid.UUID    `focus:"-"`
	VideoPreviewBlur   *entity.Media `focus:"title:Видео превью с блюром;media;hidden:list"`
	VideoPreviewBlurId *uuid.UUID    `focus:"-"`
}

func (VideoCard) TableName() string {
	return "video_cards"
}

type HtmlCard struct {
	ID   uuid.UUID `focus:"primaryKey;code:id;column:id;title:ID" validate:"required,notBlank"`
	Html string    `focus:"title:Название;filterable" validate:"required"`
}

func (HtmlCard) TableName() string {
	return "html_cards"
}

type PhotoCard struct {
	ID        uuid.UUID     `focus:"primaryKey;code:id;column:id;title:ID" validate:"required,notBlank"`
	Picture   *entity.Media `focus:"title:Фото;media;hidden:list"`
	PictureId *uuid.UUID    `focus:"-"`
}

func (PhotoCard) TableName() string {
	return "photo_cards"
}

type FormCard struct {
	ID            uuid.UUID       `focus:"primaryKey;code:id;column:id;title:ID" validate:"required,notBlank"`
	Form          *entity2.Form   `focus:"title:форма;view:select;hidden:list" validate:"omitempty,structonly"`
	FormId        *uuid.UUID      `focus:"-"`
	User          *User           `focus:"title:Спикер;view:select;viewExtra:userSelect;hidden:list" validate:"omitempty,structonly"`
	UserId        *uuid.UUID      `focus:"-"`
	FormCardsTags []FormCardsTags `focus:"-" gorm:"foreignKey:FormCardID" validate:"omitempty,unique=ID,dive,notBlank,structonly"`
	LearnMoreUrl  *string         `focus:"title:Ссылка узнать больше"`
}

func (FormCard) TableName() string {
	return "form_cards"
}

type FormCardsTags struct {
	FormCardID uuid.UUID `focus:"primaryKey;code:formCardId;column:form_card_id;title:ID" validate:"required,notBlank"`
	TagID      uuid.UUID `focus:"primaryKey;code:tagId;column:tag_id;title:ID" validate:"required,notBlank"`
	Tag        Tag       `focus:"title:Тег;view:select;hidden:list" gorm:"foreignKey:TagID"`
}

func (FormCardsTags) TableName() string {
	return "form_cards_tags"
}
