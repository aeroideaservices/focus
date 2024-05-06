package actions

import (
	"github.com/aeroideaservices/focus/media/plugin/entity"
	"github.com/google/uuid"
	"io"
)

type PageDto struct {
	ID             uuid.UUID          `json:"id"`
	Name           string             `json:"name"`
	Code           string             `json:"code"`
	Title          string             `json:"title"`
	TitleSeo       string             `json:"titleSeo"`
	Keywords       string             `json:"keywords"`
	Description    string             `json:"description"`
	DescriptionSeo string             `json:"descriptionSeo"`
	IsPublished    bool               `json:"isPublished"`
	Sort           int                `json:"sort"`
	Galleries      []GalleryInPageDto `json:"galleries"`
	//Header      *Header           `focus:"title:Header;view:select;hidden:list" validate:"omitempty,structonly"`
	//HeaderId *uuid.UUID `focus:"-" validate:"-"`
	//Footer      *Footer           `focus:"title:Footer;view:select;hidden:list" validate:"omitempty,structonly"`
	//FooterId *uuid.UUID `focus:"-" validate:"-"`

	OgType string `json:"ogType"`
}

type GetPageListResponse struct {
	PageList []PageShort `json:"pageList"`
	//PageTotal int         `json:"pageTotal"`
}
type GetPageDto struct {
	ID string `json:"id" validate:"required,sluggable"`
}

type PageShort struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
	Sort int       `json:"sort"`
}

type CreatePageRequest struct {
	Name           string                       `json:"name" validate:"required"`
	Code           string                       `json:"code" validate:"required"`
	Title          string                       `json:"title" validate:"required"`
	Description    string                       `json:"description"`
	TitleSeo       string                       `json:"titleSeo"`
	DescriptionSeo string                       `json:"descriptionSeo"`
	Keywords       string                       `json:"keywords"`
	IsPublished    *bool                        `json:"isPublished" validate:"required"`
	Sort           int                          `json:"sort"`
	Galleries      []CreateGalleryInPageRequest `json:"galleries"`

	OgType string `json:"ogType"`
}

type CreatePageResponse struct {
	ID uuid.UUID `json:"id"`
}

type PatchPageRequest struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Code           string    `json:"code"`
	Title          string    `json:"title"`
	TitleSeo       string    `json:"titleSeo"`
	Keywords       string    `json:"keywords"`
	Description    string    `json:"description"`
	DescriptionSeo string    `json:"descriptionSeo"`
	IsPublished    *bool     `json:"isPublished"`
	Sort           int       `json:"sort"`

	OgImageId *uuid.UUID `json:"ogImageId"`
	OgType    string     `json:"ogType"`
}

type GetGalleryRequest struct {
	ID string `json:"id" validate:"required,sluggable"`
}

type GalleryDto struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
	//Position    int       `json:"position"`
	IsPublished bool      `json:"isPublished"`
	Hidden      bool      `json:"hiddenMenu"`
	Cards       []CardDto `json:"cards"`
}

type GalleryInPageDto struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Position    int       `json:"position"`
	IsPublished bool      `json:"isPublished"`
	Hidden      bool      `json:"hiddenMenu"`
	Cards       []CardDto `json:"cards"`
}

type GalleryDtoWithCardsTotal struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Position    int       `json:"position"`
	IsPublished bool      `json:"isPublished"`
	Hidden      bool      `json:"hiddenMenu"`
	CardsTotal  int       `json:"cardsTotal"`
	Cards       []CardDto `json:"cards"`
}

type SearchGalleryRequest struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Code       string    `json:"code"`
	Position   int       `json:"position"`
	CardsTotal int       `json:"cardsTotal"`
	Cards      []CardDto `json:"cards"`
	//Hidden     bool
	//IsPublished bool
}
type PatchGalleryPosition struct {
	PageID      uuid.UUID `json:"pageId" validate:"required"`
	GalleryID   uuid.UUID `json:"galleryId" validate:"required"`
	OldPosition *int      `json:"oldPosition" validate:"required"`
	NewPosition *int      `json:"newPosition" validate:"required"`
}

type UpdateGalleryRequest struct {
	ID          uuid.UUID           `json:"id"`
	Name        string              `json:"name"`
	Code        string              `json:"code"`
	Position    int                 `json:"position"`
	CardsTotal  int                 `json:"cardsTotal"`
	Cards       []CreateCardRequest `json:"cards"`
	Hidden      *bool               `json:"hiddenMenu"`
	IsPublished *bool               `json:"isPublished"`
}

type PatchGalleryNameRequest struct {
	ID   uuid.UUID `json:"id" validate:"required"`
	Name string    `json:"name" validate:"required"`
}

type PatchCardPosition struct {
	GalleryID   uuid.UUID `json:"galleryId" validate:"required"`
	CardID      uuid.UUID `json:"cardId" validate:"required"`
	OldPosition *int      `json:"oldPosition" validate:"required"`
	NewPosition *int      `json:"newPosition" validate:"required"`
}

type CreateGalleryInPageRequest struct {
	Name        string              `json:"name"`
	Code        string              `json:"code"`
	Position    int                 `json:"position"`
	CardsTotal  int                 `json:"cardsTotal"`
	Cards       []CreateCardRequest `json:"cards"`
	Hidden      bool                `json:"hiddenMenu"`
	IsPublished bool                `json:"isPublished"`
}

type CreateGalleryRequest struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
	//Position    int                 `json:"position"`
	CardsTotal  int                 `json:"cardsTotal"`
	Cards       []CreateCardRequest `json:"cards"`
	Hidden      *bool               `json:"hiddenMenu" validate:"required"`
	IsPublished *bool               `json:"isPublished" validate:"required"`
}

type CreateGalleryResponse struct {
	ID uuid.UUID `json:"id"`
}

type CardDto struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Code        string          `json:"code"`
	Type        string          `json:"type"`
	Position    int             `json:"position"`
	IsPublished bool            `json:"isPublished"`
	RegularCard *RegularCardDto `json:"regularCard"`
	VideoCard   *VideoCardDto   `json:"videoCard"`
	HtmlCard    *HtmlCardDto    `json:"htmlCard"`
	PhotoCard   *PhotoCardDto   `json:"photoCard"`
	FormCard    *FormCardDto    `json:"formCard"`

	Title       string `json:"title"`
	Description string `json:"description"`
	OgType      string `json:"ogType"`
}

type CardDtoWithoutPosition struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Code        string          `json:"code"`
	Type        string          `json:"type"`
	IsPublished bool            `json:"isPublished"`
	RegularCard *RegularCardDto `json:"regularCard"`
	VideoCard   *VideoCardDto   `json:"videoCard"`
	HtmlCard    *HtmlCardDto    `json:"htmlCard"`
	PhotoCard   *PhotoCardDto   `json:"photoCard"`
	FormCard    *FormCardDto    `json:"formCard"`

	Title       string `json:"title"`
	Description string `json:"description"`
	OgType      string `json:"ogType"`
}

type RegularCardDto struct {
	PreviewText      string        `json:"previewText"`
	DetailText       string        `json:"detailText"`
	Inverted         bool          `json:"inverted"`
	Video            *entity.Media `json:"video"`
	VideoLite        *entity.Media `json:"videoLite"`
	VideoPreview     *entity.Media `json:"videoPreview"`
	VideoPreviewBlur *entity.Media `json:"videoPreviewBlur"`
	User             *UserDto      `json:"user"`
	Tags             []TagDto      `json:"tags"`
	LearnMoreUrl     *string       `json:"learnMoreUrl"`
}

type VideoCardDto struct {
	Video            *entity.Media `json:"video"`
	VideoLite        *entity.Media `json:"videoLite"`
	VideoPreview     *entity.Media `json:"videoPreview"`
	VideoPreviewBlur *entity.Media `json:"videoPreviewBlur"`
}

type HtmlCardDto struct {
	Html string `json:"html"`
}

type PhotoCardDto struct {
	Picture *entity.Media `json:"picture"`
}

type FormCardDto struct {
	Form         *FormDto      `json:"form"`
	Background   *entity.Media `json:"background"`
	User         *UserDto      `json:"user"`
	Tags         []TagDto      `json:"tags"`
	LearnMoreUrl *string       `json:"learnMoreUrl"`
}

type FormDto struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CreateCardRequest struct {
	GalleryIds  []uuid.UUID               `json:"galleryIds"`
	Type        string                    `json:"type"  validate:"required"`
	Name        string                    `json:"name"`
	Code        string                    `json:"code"`
	Position    int                       `json:"position"`
	IsPublished *bool                     `json:"isPublished"  validate:"required"`
	RegularCard *CreateRegularCardRequest `json:"regularCard"`
	VideoCard   *CreateVideoCardRequest   `json:"videoCard"`
	HtmlCard    *CreateHtmlCardRequest    `json:"htmlCard"`
	PhotoCard   *CreatePhotoCardRequest   `json:"photoCard"`
	FormCard    *CreateFormCardRequest    `json:"formCard"`

	Title       string `json:"title"`
	Description string `json:"description"`
	OgType      string `json:"ogType"`
}

type UpdateCardRequest struct {
	ID         uuid.UUID   `json:"id"`
	Name       string      `json:"name"`
	Code       string      `json:"code"`
	GalleryIds []uuid.UUID `json:"galleryIds"`
	Type       string      `json:"type"`
	//Position    int                       `json:"position"`
	IsPublished *bool                     `json:"isPublished"`
	RegularCard *CreateRegularCardRequest `json:"regularCard"`
	VideoCard   *CreateVideoCardRequest   `json:"videoCard"`
	HtmlCard    *CreateHtmlCardRequest    `json:"htmlCard"`
	PhotoCard   *CreatePhotoCardRequest   `json:"photoCard"`
	FormCard    *CreateFormCardRequest    `json:"formCard"`

	Title       string `json:"title"`
	Description string `json:"description"`
	OgType      string `json:"ogType"`
}

type CreateCardResponse struct {
	ID uuid.UUID `json:"id"`
}

type CreateTagResponse struct {
	ID uuid.UUID `json:"id"`
}

type GetCardRequest struct {
	ID string `json:"id" validate:"required,sluggable"`
}

type PatchUserRequest struct {
	CardId uuid.UUID `json:"cardId" validate:"required"`
	UserId uuid.UUID `json:"userId" validate:"required"`
}

type PatchPreviewTextRequest struct {
	CardId      uuid.UUID `json:"cardId" validate:"required"`
	PreviewText string    `json:"previewText" validate:"required"`
}

type PatchDetailTextRequest struct {
	CardId     uuid.UUID `json:"cardId" validate:"required"`
	DetailText string    `json:"detailText" validate:"required"`
}

type PatchLearnMoreUrlRequest struct {
	CardId       uuid.UUID `json:"cardId" validate:"required"`
	LearnMoreUrl *string   `json:"learnMoreUrl"`
}

type CreateRegularCardRequest struct {
	PreviewText        string      `json:"previewText"`
	DetailText         string      `json:"detailText"`
	Inverted           *bool       `json:"inverted"`
	VideoId            *uuid.UUID  `json:"videoId" validate:"required"`
	VideoLiteId        *uuid.UUID  `json:"videoLiteId" validate:"required"`
	VideoPreviewId     *uuid.UUID  `json:"videoPreviewId" validate:"required"`
	VideoPreviewBlurId *uuid.UUID  `json:"videoPreviewBlurId" validate:"required"`
	Tags               []uuid.UUID `json:"tags"`
	UserId             *uuid.UUID  `json:"userId"`
	LearnMoreUrl       *string     `json:"learnMoreUrl"`
}

type CreateVideoCardRequest struct {
	VideoId            *uuid.UUID `json:"videoId"`
	VideoLiteId        *uuid.UUID `json:"videoLiteId"`
	VideoPreviewId     *uuid.UUID `json:"videoPreviewId"`
	VideoPreviewBlurId *uuid.UUID `json:"videoPreviewBlurId"`
}

type CreateHtmlCardRequest struct {
	Html string `json:"html"`
}

type CreatePhotoCardRequest struct {
	PictureId *uuid.UUID `json:"pictureId"`
}

type CreateFormCardRequest struct {
	FormId       *uuid.UUID  `json:"formId"`
	BackgroundId *uuid.UUID  `json:"backgroundId"`
	UserId       *uuid.UUID  `json:"userId"`
	Tags         []uuid.UUID `json:"tags"`
	LearnMoreUrl *string     `json:"learnMoreUrl"`
}

type TagDto struct {
	ID   uuid.UUID `json:"id"`
	Text string    `json:"text" validate:"required"`
	Link string    `json:"link" validate:"required"`
}

type TagSearchDto struct {
	ID      uuid.UUID   `json:"id"`
	CardIds []uuid.UUID `json:"cardIds"`
	Text    string      `json:"text" validate:"required"`
	Link    string      `json:"link" validate:"required"`
}

type TagDtoRequest struct {
	ID   uuid.UUID `json:"id"`
	Text string    `json:"text"`
	Link string    `json:"link"`
}

type UserDto struct {
	ID        uuid.UUID     `json:"id"`
	FirstName string        `json:"firstName"`
	LastName  string        `json:"lastName"`
	Position  string        `json:"position"`
	Picture   *entity.Media `json:"picture"`
}

type CreateVideoRequest struct {
	Filename string        `validate:"required,min=3"`
	Size     int64         `validate:""`
	FolderId *uuid.UUID    `validate:"omitempty,notBlank"`
	File     io.ReadSeeker `validate:"required"`
}

type CreateVideoResponse struct {
	VideoId          uuid.UUID `json:"videoId"`
	VideoLiteId      uuid.UUID `json:"videoLiteId"`
	PreviewId        uuid.UUID `json:"videoPreviewId"`
	PreviewBlurredId uuid.UUID `json:"videoPreviewBlurId"`
}
