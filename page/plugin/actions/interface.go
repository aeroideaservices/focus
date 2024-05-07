package actions

import (
	"context"
	"github.com/aeroideaservices/focus/page/plugin/entity"
	"github.com/google/uuid"
)

type PageRepository interface {
	GetById(ctx context.Context, id uuid.UUID) (*entity.Page, error)
	GetByIdWithoutAssociate(ctx context.Context, id uuid.UUID) (*entity.Page, error)
	GetList(ctx context.Context, fields []string, allowInactive bool) ([]PageShort, error)
	Create(ctx context.Context, page *entity.Page) (*uuid.UUID, error)
	Delete(ctx context.Context, pageID uuid.UUID) error
	PatchProperties(ctx context.Context, page *entity.Page) error
	PatchGalleryPosition(ctx context.Context, dto *PatchGalleryPosition) error
	GetLastPosition(ctx context.Context, pageID uuid.UUID) (int, error)
	CreatePagesGalleries(ctx context.Context, pagesGalleries []entity.PagesGalleries) (error, bool)
	DeletePagesGalleries(ctx context.Context, pageID uuid.UUID, galleryIDs []uuid.UUID) error
	UpdatePublish(ctx context.Context, pageID uuid.UUID, publish *bool) error
}

type GalleryRepository interface {
	GetByCode(ctx context.Context, code string) (*entity.Gallery, error)
	GetById(ctx context.Context, id uuid.UUID) (*entity.Gallery, error)
	GetByIdWithoutAssociate(ctx context.Context, id uuid.UUID) (*entity.Gallery, error)
	GetListByPageId(ctx context.Context, pageId uuid.UUID) ([]entity.Gallery, error)
	Create(ctx context.Context, gallery *entity.Gallery) (*uuid.UUID, error)
	GetListWithSearch(ctx context.Context, searchValue string) ([]entity.Gallery, error)
	Update(ctx context.Context, gallery *entity.Gallery) error
	PatchName(ctx context.Context, gallery *entity.Gallery) error
	PatchCardPosition(ctx context.Context, dto *PatchCardPosition) error
	DeleteList(ctx context.Context, galleryIds []uuid.UUID) error
	GetLastPosition(ctx context.Context, galleryID uuid.UUID) (int, error)
	CreateGalleriesCards(ctx context.Context, galleriesCards []entity.GalleriesCards) (error, bool)
	DeleteGalleriesCards(ctx context.Context, galleryID uuid.UUID, cardIDs []uuid.UUID) error
	UpdateHidden(ctx context.Context, galleryID uuid.UUID, hidden *bool) error
	UpdatePublish(ctx context.Context, galleryID uuid.UUID, publish *bool) error
}

type CardRepository interface {
	GetById(ctx context.Context, cardId uuid.UUID) (*entity.Card, error)
	GetByIdWithoutAssociate(ctx context.Context, cardId uuid.UUID) (*entity.Card, error)
	GetListByGalleryId(ctx context.Context, galleryId uuid.UUID) ([]entity.Card, error)
	Create(ctx context.Context, card *entity.Card, galleriesCards []entity.GalleriesCards) (*uuid.UUID, error)
	GetListWithSearch(ctx context.Context, searchValue string, name string) ([]entity.Card, error)
	Update(ctx context.Context, card *entity.Card, galleriesCards []entity.GalleriesCards) error
	Delete(ctx context.Context, cardId uuid.UUID) error
	PatchUser(ctx context.Context, card *entity.Card) error
	PatchPreviewText(ctx context.Context, card *entity.Card) error
	PatchDetailText(ctx context.Context, card *entity.Card) error
	PatchLearnMoreUrl(ctx context.Context, card *entity.Card) error
	PatchTags(ctx context.Context, cardId uuid.UUID, tags []entity.Tag) error
	CreateRegularTags(ctx context.Context, cardId uuid.UUID, tagIds []uuid.UUID) (error, bool)
	DeleteRegularTags(ctx context.Context, cardId uuid.UUID, tagIds []uuid.UUID) error
	UpdatePublish(ctx context.Context, cardID uuid.UUID, publish *bool) error
	GetLastPositionInGalley(ctx context.Context, galleryID uuid.UUID) (int, error)
	UpdateInverted(ctx context.Context, cardID uuid.UUID, inverted *bool) error
}

type TagRepository interface {
	//GetListByCardId(ctx context.Context, regularCardId uuid.UUID) ([]entity.ExternalTag, []entity.InternalTag, error)
	UpdateIsDetailLink(ctx context.Context, tagID uuid.UUID, value *bool) error
	GetListWithSearch(ctx context.Context, searchValue string) ([]entity.Tag, map[uuid.UUID][]uuid.UUID, error)
	GetById(ctx context.Context, tagId uuid.UUID) (*entity.Tag, error)
	Create(ctx context.Context, tag *entity.Tag) (*uuid.UUID, error)
	Update(ctx context.Context, tag *entity.Tag) error
	Delete(ctx context.Context, tagId uuid.UUID) error
}

type UserRepository interface {
	GetById(ctx context.Context, userId uuid.UUID) (*entity.User, error)
}

type CopierInterface interface {
	Copy(toValue interface{}, fromValue interface{}) (err error)
}
