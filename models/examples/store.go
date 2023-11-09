package examples

import (
	"github.com/aeroideaservices/focus/media/plugin/entity"
	"github.com/google/uuid"
)

type Store struct {
	ID           uuid.UUID     `focus:"title:ID;primaryKey;code:id;column:id"`
	Name         string        `focus:"title:Название;filterable" validate:"required,min=3,max=50"`
	Latitude     float64       `focus:"title:Широта" validate:"required,lt=100"`
	Longitude    float64       `focus:"title:Долгота" validate:"required,lt=100"`
	ContactEmail *string       `focus:"title:Контактный email;block:Контакты;view:emailInput" validate:"omitempty,pattern=email"`
	ContactPhone *string       `focus:"title:Контактный телефон;block:Контакты;view:phoneInput" validate:"omitempty,pattern=phone"`
	OpeningTime  string        `focus:"title:Время открытия;view:timePickerInput" validate:"required,datetime=15:04:05"`
	ClosingTime  string        `focus:"title:Время закрытия;view:timePickerInput" validate:"required,datetime=15:04:05"`
	Products     []Product     `focus:"title:Товары;many2many:stores_products;viewExtra:selectProducts;joinSort:sort" gorm:"many2many:stores_products" validate:"omitempty,unique=ID,dive,notBlank,structonly"`
	Image        *entity.Media `focus:"title:Изображение;media;viewExtra:storesMedia" validate:"omitempty"`
	ImageId      *uuid.UUID    `focus:"-"`
	Description  *string       `focus:"title:Описание магазина;view:editorJs;viewExtra:storeEditorJs" validate:"omitempty,json"`
}

func (Store) TableName() string {
	return "stores"
}

func (Store) ModelTitle() string {
	return "Магазины"
}
