package repositories

import (
	"context"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"pages/pkg/page/plugin/actions"
	"pages/pkg/page/plugin/entity"
	"sort"
)

type PageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) *PageRepository {
	return &PageRepository{
		db: db,
	}
}

func (r *PageRepository) GetById(ctx context.Context, id uuid.UUID) (
	*entity.Page, error,
) {
	page := &entity.Page{}
	db := r.db.WithContext(ctx).Model(entity.Page{}).
		Preload(
			"PagesGalleries", func(db *gorm.DB) *gorm.DB {
				return db.Order("pages_galleries.position ASC")
			},
		).
		Preload(
			"PagesGalleries.Gallery.GalleriesCards", func(db *gorm.DB) *gorm.DB {
				return db.Order("galleries_cards.position ASC")
			},
		).
		Preload("PagesGalleries.Gallery.GalleriesCards.Card").
		Preload("PagesGalleries.Gallery.GalleriesCards.Card.HtmlCard").
		Preload("PagesGalleries.Gallery.GalleriesCards.Card.VideoCard").Preload("PagesGalleries.Gallery.GalleriesCards.Card.VideoCard.Video").Preload("PagesGalleries.Gallery.GalleriesCards.Card.VideoCard.VideoLite").
		Preload("PagesGalleries.Gallery.GalleriesCards.Card.VideoCard.VideoPreview").Preload("PagesGalleries.Gallery.GalleriesCards.Card.VideoCard.VideoPreviewBlur").
		Preload("PagesGalleries.Gallery.GalleriesCards.Card.RegularCard.Video").Preload("PagesGalleries.Gallery.GalleriesCards.Card.RegularCard.VideoLite").Preload("PagesGalleries.Gallery.GalleriesCards.Card.RegularCard.VideoPreview").
		Preload("PagesGalleries.Gallery.GalleriesCards.Card.RegularCard.VideoPreviewBlur").Preload("PagesGalleries.Gallery.GalleriesCards.Card.RegularCard.User").
		Preload("PagesGalleries.Gallery.GalleriesCards.Card.RegularCard.RegularCardsTags.Tag").
		Preload("PagesGalleries.Gallery.GalleriesCards.Card.RegularCard.User.Picture").
		Preload("PagesGalleries.Gallery.GalleriesCards.Card.PhotoCard").
		Preload("PagesGalleries.Gallery.GalleriesCards.Card.PhotoCard.Picture").
		Preload("PagesGalleries.Gallery.GalleriesCards.Card.FormCard").Preload("PagesGalleries.Gallery.GalleriesCards.Card.FormCard.Background").
		Preload("PagesGalleries.Gallery.GalleriesCards.Card.FormCard.User").Preload("PagesGalleries.Gallery.GalleriesCards.Card.FormCard.User.Picture").
		Preload("PagesGalleries.Gallery.GalleriesCards.Card.FormCard.FormCardsTags.Tag").Preload("PagesGalleries.Gallery.GalleriesCards.Card.FormCard.Form").
		Where("id", id)
	err := db.First(page).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrapf(err, "pages with id %s not found", id)
	}
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (r *PageRepository) GetByIdWithoutAssociate(ctx context.Context, id uuid.UUID) (*entity.Page, error) {
	page := &entity.Page{}
	db := r.db.WithContext(ctx).Omit(clause.Associations).Model(entity.Page{}).Where("id", id)
	err := db.First(page).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrapf(err, "pages with id %s not found", id)
	}
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (r *PageRepository) GetList(
	ctx context.Context, fields []string, allowInactive bool,
) (list []actions.PageShort, err error) {
	//db := r.db.WithContext(ctx).Model(page_entity.PageShort{})
	db := r.db.WithContext(ctx).Select(fields).Order("pages.sort ASC").Model(entity.Page{})

	if !allowInactive {
		db = db.Scopes(pageIsPublished)
	}

	err = db.Scan(&list).Error

	return
}

//func (r *PageRepository) Create(ctx context.Context, page *entity.Page) (*uuid.UUID, error) {
//	tx := r.db.WithContext(ctx).Begin()
//
//	if err := tx.Omit(clause.Associations).Create(page).Error; err != nil {
//		tx.Rollback()
//		return nil, err
//	}
//
//	for _, pagesGallery := range page.PagesGalleries {
//		if err := tx.Omit(clause.Associations).Create(&pagesGallery.Gallery).Error; err != nil {
//			tx.Rollback()
//			return nil, err
//		}
//		if err := tx.Omit(clause.Associations).Create(&pagesGallery).Error; err != nil {
//			tx.Rollback()
//			return nil, err
//		}
//
//		for _, galleriesCard := range pagesGallery.Gallery.GalleriesCards {
//
//			if galleriesCard.Card.RegularCard != nil {
//				if err := tx.Omit(clause.Associations).Create(galleriesCard.Card.RegularCard).Error; err != nil {
//					tx.Rollback()
//					return nil, err
//				}
//				for _, tag := range galleriesCard.Card.RegularCard.RegularCardsTags {
//					if err := tx.Omit(clause.Associations).Create(tag).Error; err != nil {
//						tx.Rollback()
//						return nil, err
//					}
//				}
//			}
//
//			if galleriesCard.Card.VideoCard != nil {
//				if err := tx.Omit(clause.Associations).Create(galleriesCard.Card.VideoCard).Error; err != nil {
//					tx.Rollback()
//					return nil, err
//				}
//			}
//
//			if galleriesCard.Card.HtmlCard != nil {
//				if err := tx.Omit(clause.Associations).Omit(clause.Associations).Create(galleriesCard.Card.HtmlCard).Error; err != nil {
//					tx.Rollback()
//					return nil, err
//				}
//			}
//
//			if galleriesCard.Card.PhotoCard != nil {
//				if err := tx.Omit(clause.Associations).Omit(clause.Associations).Create(galleriesCard.Card.PhotoCard).Error; err != nil {
//					tx.Rollback()
//					return nil, err
//				}
//			}
//
//			if galleriesCard.Card.FormCard != nil {
//				if err := tx.Omit(clause.Associations).Omit(clause.Associations).Create(galleriesCard.Card.FormCard).Error; err != nil {
//					tx.Rollback()
//					return nil, err
//				}
//				for _, tag := range galleriesCard.Card.FormCard.FormCardsTags {
//
//					if err := tx.Omit(clause.Associations).Create(tag).Error; err != nil {
//						tx.Rollback()
//						return nil, err
//					}
//
//				}
//			}
//
//			if err := tx.Omit(clause.Associations).Create(&galleriesCard.Card).Error; err != nil {
//				tx.Rollback()
//				return nil, err
//			}
//
//			if err := tx.Omit(clause.Associations).Create(&galleriesCard).Error; err != nil {
//				tx.Rollback()
//				return nil, err
//			}
//
//		}
//	}
//
//	if err := tx.Commit().Error; err != nil {
//		return nil, err
//	}
//
//	return &page.ID, nil
//}

func (r *PageRepository) Create(ctx context.Context, page *entity.Page) (*uuid.UUID, error) {
	tx := r.db.WithContext(ctx).Begin()

	if err := createPage(tx, page); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := createPagesGalleries(tx, page); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &page.ID, nil
}

func createPage(tx *gorm.DB, page *entity.Page) error {
	return tx.Omit(clause.Associations).Create(page).Error
}

func createPagesGalleries(tx *gorm.DB, page *entity.Page) error {
	for _, pagesGallery := range page.PagesGalleries {
		if err := createGallery(tx, &pagesGallery.Gallery); err != nil {
			return err
		}
		if err := createPagesGallery(tx, &pagesGallery); err != nil {
			return err
		}
		if err := createGalleriesCards(tx, &pagesGallery.Gallery); err != nil {
			return err
		}
	}
	return nil
}

func createGallery(tx *gorm.DB, gallery *entity.Gallery) error {
	return tx.Omit(clause.Associations).Create(gallery).Error
}

func createPagesGallery(tx *gorm.DB, pagesGallery *entity.PagesGalleries) error {
	return tx.Omit(clause.Associations).Create(pagesGallery).Error
}

func createGalleriesCards(tx *gorm.DB, gallery *entity.Gallery) error {
	for _, galleriesCard := range gallery.GalleriesCards {
		if err := createCard(tx, galleriesCard.Card); err != nil {
			return err
		}
		if err := createGalleriesCard(tx, &galleriesCard); err != nil {
			return err
		}
	}
	return nil
}

func createGalleriesCard(tx *gorm.DB, galleriesCard *entity.GalleriesCards) error {
	return tx.Omit(clause.Associations).Create(galleriesCard).Error
}

func createCard(tx *gorm.DB, card *entity.Card) error {
	switch {
	case card.RegularCard != nil:
		if err := createRegularCard(tx, card.RegularCard); err != nil {
			return err
		}
	case card.VideoCard != nil:
		if err := createVideoCard(tx, card.VideoCard); err != nil {
			return err
		}
	case card.HtmlCard != nil:
		if err := createHtmlCard(tx, card.HtmlCard); err != nil {
			return err
		}
	case card.PhotoCard != nil:
		if err := createPhotoCard(tx, card.PhotoCard); err != nil {
			return err
		}
	case card.FormCard != nil:
		if err := createFormCard(tx, card.FormCard); err != nil {
			return err
		}
	}
	return tx.Omit(clause.Associations).Create(card).Error
}

func createRegularCard(tx *gorm.DB, regularCard *entity.RegularCard) error {
	if err := tx.Omit(clause.Associations).Create(regularCard).Error; err != nil {
		return err
	}

	for _, tag := range regularCard.RegularCardsTags {
		if err := tx.Omit(clause.Associations).Create(tag).Error; err != nil {
			return err
		}
	}

	return nil
}

func createVideoCard(tx *gorm.DB, videoCard *entity.VideoCard) error {
	if err := tx.Omit(clause.Associations).Create(videoCard).Error; err != nil {
		return err
	}
	return nil
}

func createHtmlCard(tx *gorm.DB, htmlCard *entity.HtmlCard) error {
	if err := tx.Omit(clause.Associations).Create(htmlCard).Error; err != nil {
		return err
	}
	return nil
}

func createPhotoCard(tx *gorm.DB, photoCard *entity.PhotoCard) error {
	if err := tx.Omit(clause.Associations).Create(photoCard).Error; err != nil {
		return err
	}
	return nil
}

func createFormCard(tx *gorm.DB, formCard *entity.FormCard) error {
	if err := tx.Omit(clause.Associations).Create(formCard).Error; err != nil {
		return err
	}

	for _, tag := range formCard.FormCardsTags {
		if err := tx.Omit(clause.Associations).Create(tag).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *PageRepository) Delete(ctx context.Context, pageId uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("id = ?", pageId).Delete(entity.Page{}).Error
	return err
}

func (r *PageRepository) PatchProperties(ctx context.Context, page *entity.Page) error {
	err := r.db.WithContext(ctx).Model(page).Where("id = ?", page.ID).
		Updates(
			page,
		).Error

	return err
}
func (r *PageRepository) UpdatePublish(ctx context.Context, pageId uuid.UUID, publish *bool) error {
	err := r.db.WithContext(ctx).Model(&entity.Page{}).Where("id = ?", pageId).
		Updates(
			map[string]interface{}{
				"is_published": publish,
			},
		).Error
	return err
}

func (r *PageRepository) PatchGalleryPosition(ctx context.Context, dto *actions.PatchGalleryPosition) error {
	tx := r.db.WithContext(ctx).Begin()

	err2 := r.UpdateGalleriesPositionAfterPatch(ctx, dto, tx)
	if err2 != nil {
		return err2
	}

	err := tx.WithContext(ctx).Table("pages_galleries").
		Where("gallery_id = ?", dto.GalleryID).
		Where("pages_id = ?", dto.PageID).
		Update("position", dto.NewPosition).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return err
}

func (r *PageRepository) UpdateGalleriesPositionAfterPatch(
	ctx context.Context, dto *actions.PatchGalleryPosition, tx *gorm.DB,
) error {
	if *dto.OldPosition > *dto.NewPosition {
		err := tx.WithContext(ctx).Table("pages_galleries").Where("pages_id = ?", dto.PageID).Where(
			"position >= ? and position < ?", dto.NewPosition, dto.OldPosition,
		).Update("position", gorm.Expr("position + 1")).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if *dto.OldPosition < *dto.NewPosition {
		err := tx.WithContext(ctx).Table("pages_galleries").Where("pages_id = ?", dto.PageID).Where(
			"position <= ? and position > ?", dto.NewPosition, dto.OldPosition,
		).Update("position", gorm.Expr("position - 1")).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

func pageIsPublished(db *gorm.DB) *gorm.DB {
	//published := clause.Column{Table: page_entity.Page{}.TableName(), Name: "is_published"}
	return db.
		Where("is_published IS TRUE")
}

func (r *PageRepository) CreatePagesGalleries(ctx context.Context, pagesGalleries []entity.PagesGalleries) (
	error, bool,
) {
	err := r.db.WithContext(ctx).Omit(clause.Associations).Create(&pagesGalleries).Error

	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == UniqueViolationErr {
			return errors.BadRequest.Wrapf(pgErr, "gallery already linked to the page, detail: %s", pgErr.Detail), true
		}
	}

	return err, false
}

func (r *PageRepository) DeletePagesGalleries(ctx context.Context, pageID uuid.UUID, galleryIDs []uuid.UUID) error {
	for _, galleryId := range galleryIDs {
		ok := r.CheckLink(ctx, pageID, galleryId)
		if !ok {
			return errors.NotFound.Newf("gallery with id: %s is not linked to the page", galleryId.String())
		}
	}

	var pagesGalleries []entity.PagesGalleries
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).
		Where("pages_id = ?", pageID).
		Where("gallery_id in ?", galleryIDs).
		Delete(&pagesGalleries).Error
	if err != nil {
		return err
	}

	err = r.UpdateGalleriesPositionAfterDelete(ctx, pageID, pagesGalleries, err)
	return err
}

func (r *PageRepository) UpdateGalleriesPositionAfterDelete(
	ctx context.Context, pageID uuid.UUID, pagesGalleries []entity.PagesGalleries, err error,
) error {
	sort.Slice(
		pagesGalleries, func(i, j int) bool {
			return pagesGalleries[i].Position > pagesGalleries[j].Position
		},
	)

	for _, gallery := range pagesGalleries {
		err = r.db.WithContext(ctx).Table("pages_galleries").Where("pages_id = ?", pageID).Where(
			"position > ?", gallery.Position,
		).Update("position", gorm.Expr("position - 1")).Error
	}
	return err
}

func (r *PageRepository) GetLastPosition(ctx context.Context, pageID uuid.UUID) (int, error) {
	var result int
	var count int64
	r.db.WithContext(ctx).Table("pages_galleries").Where(
		"pages_id = ?", pageID,
	).Count(&count)
	if count == 0 {
		return 0, nil
	}
	err := r.db.WithContext(ctx).Table("pages_galleries").Where(
		"pages_id = ?", pageID,
	).Select("max(position)").Row().Scan(&result)

	if errors.Is(err, gorm.ErrInvalidField) {
		return 0, nil
	}

	return result, err
}

func (r *PageRepository) CheckLink(ctx context.Context, pageID uuid.UUID, galleryID uuid.UUID) bool {
	pagesGalleries := &entity.PagesGalleries{}
	err := r.db.WithContext(ctx).Omit(clause.Associations).Model(entity.PagesGalleries{}).
		Where("pages_id = ?", pageID).
		Where("gallery_id = ?", galleryID).
		First(pagesGalleries).Error
	if err != nil || pagesGalleries == nil {
		return false
	}
	return true
}
