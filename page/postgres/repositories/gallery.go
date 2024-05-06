package repositories

import (
	"context"
	media_entity "github.com/aeroideaservices/focus/media/plugin/entity"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"pages/pkg/page/plugin/actions"
	"pages/pkg/page/plugin/entity"
	"sort"
)

type GalleryRepository struct {
	db *gorm.DB
}

func NewGalleryRepository(db *gorm.DB) *GalleryRepository {
	return &GalleryRepository{
		db: db,
	}
}

func (r *GalleryRepository) GetByCode(ctx context.Context, code string) (*entity.Gallery, error) {
	gallery := &entity.Gallery{}

	db := r.db.WithContext(ctx).Model(entity.Gallery{}).Where("code", code)

	err := db.First(gallery).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrapf(err, "gallery with code %s not found", code)
	}
	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (r *GalleryRepository) GetById(ctx context.Context, id uuid.UUID) (*entity.Gallery, error) {
	gallery := &entity.Gallery{}

	db := r.db.WithContext(ctx).Model(entity.Gallery{}).
		Preload(
			"GalleriesCards", func(db *gorm.DB) *gorm.DB {
				return db.Order("galleries_cards.position ASC")
			},
		).
		Preload("GalleriesCards.Card").
		Preload("GalleriesCards.Card.HtmlCard").
		Preload("GalleriesCards.Card.VideoCard").Preload("GalleriesCards.Card.VideoCard.Video").Preload("GalleriesCards.Card.VideoCard.VideoLite").Preload("GalleriesCards.Card.VideoCard.VideoPreview").Preload("GalleriesCards.Card.VideoCard.VideoPreviewBlur").
		Preload("GalleriesCards.Card.RegularCard.Video").Preload("GalleriesCards.Card.RegularCard.VideoLite").Preload("GalleriesCards.Card.RegularCard.VideoPreview").Preload("GalleriesCards.Card.RegularCard.VideoPreviewBlur").Preload("GalleriesCards.Card.RegularCard.User").
		Preload("GalleriesCards.Card.RegularCard.RegularCardsTags.Tag").
		Preload("GalleriesCards.Card.RegularCard.User.Picture").
		Preload("GalleriesCards.Card.PhotoCard").Preload("GalleriesCards.Card.PhotoCard.Picture").
		Preload("GalleriesCards.Card.FormCard").Preload("GalleriesCards.Card.FormCard.Background").Preload("GalleriesCards.Card.FormCard.User").Preload("GalleriesCards.Card.FormCard.User.Picture").
		Preload("GalleriesCards.Card.FormCard.FormCardsTags.Tag").Preload("GalleriesCards.Card.FormCard.Form").
		Where("id", id)
	err := db.First(gallery).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrapf(err, "gallery with id %s not found", id)
	}
	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (r *GalleryRepository) GetByIdWithoutAssociate(ctx context.Context, id uuid.UUID) (*entity.Gallery, error) {
	gallery := &entity.Gallery{}

	db := r.db.WithContext(ctx).Omit(clause.Associations).Model(entity.Gallery{}).
		Where("id", id)
	err := db.First(gallery).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrapf(err, "gallery with id %s not found", id)
	}
	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (r *GalleryRepository) GetListByPageId(ctx context.Context, pageId uuid.UUID) ([]entity.Gallery, error) {
	var galleries []entity.Gallery
	db := r.db.WithContext(ctx).Select("galleries.*, pages_galleries.position").Table("pages_galleries").
		Joins("left join galleries on pages_galleries.gallery_id=galleries.id").
		Where("pages_id=?", pageId).
		Find(&galleries)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, db.Error
		}
	}

	return galleries, nil
}

func (r *GalleryRepository) GetListWithSearch(ctx context.Context, searchValue string) (
	[]entity.Gallery, error,
) {
	var galleries []entity.Gallery

	db := r.db.WithContext(ctx).Model(entity.Gallery{}).
		Preload(
			"GalleriesCards", func(db *gorm.DB) *gorm.DB {
				return db.Order("galleries_cards.position ASC")
			},
		).
		Preload("GalleriesCards.Card").
		Preload("GalleriesCards.Card.HtmlCard").
		Preload("GalleriesCards.Card.VideoCard").Preload("GalleriesCards.Card.VideoCard.Video").Preload("GalleriesCards.Card.VideoCard.VideoLite").Preload("GalleriesCards.Card.VideoCard.VideoPreview").Preload("GalleriesCards.Card.VideoCard.VideoPreviewBlur").
		Preload("GalleriesCards.Card.RegularCard.Video").Preload("GalleriesCards.Card.RegularCard.VideoLite").Preload("GalleriesCards.Card.RegularCard.VideoPreview").Preload("GalleriesCards.Card.RegularCard.VideoPreviewBlur").Preload("GalleriesCards.Card.RegularCard.User").
		Preload("GalleriesCards.Card.RegularCard.RegularCardsTags.Tag").
		Preload("GalleriesCards.Card.RegularCard.User.Picture").
		Preload("GalleriesCards.Card.PhotoCard").Preload("GalleriesCards.Card.PhotoCard.Picture").
		Preload("GalleriesCards.Card.FormCard").Preload("GalleriesCards.Card.FormCard.Background").Preload("GalleriesCards.Card.FormCard.User").Preload("GalleriesCards.Card.FormCard.User.Picture").
		Preload("GalleriesCards.Card.FormCard.FormCardsTags.Tag").Preload("GalleriesCards.Card.FormCard.Form").
		Where("lower(name) LIKE lower(?)", "%"+searchValue+"%").
		Or("lower(code) LIKE lower(?)", "%"+searchValue+"%")
	err := db.Find(&galleries).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrapf(err, "galleries with searchValue %s not found", searchValue)
	}
	if err != nil {
		return nil, err
	}
	return galleries, nil
}

func (r *GalleryRepository) Create(ctx context.Context, gallery *entity.Gallery) (*uuid.UUID, error) {
	tx := r.db.WithContext(ctx).Begin()

	if err := tx.Omit(clause.Associations).Create(gallery).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, galleriesCard := range gallery.GalleriesCards {

		if galleriesCard.Card.VideoCard != nil && galleriesCard.Card.VideoCard.VideoId != nil {
			var videoMedia media_entity.Media
			if err := tx.First(&videoMedia, *galleriesCard.Card.VideoCard.VideoId).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			galleriesCard.Card.VideoCard.Video = &videoMedia
		}

		if galleriesCard.Card.RegularCard != nil && galleriesCard.Card.RegularCard.VideoId != nil {
			var regularVideoMedia media_entity.Media
			if err := tx.First(&regularVideoMedia, *galleriesCard.Card.RegularCard.VideoId).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			galleriesCard.Card.RegularCard.Video = &regularVideoMedia
		}

		if galleriesCard.Card.RegularCard != nil && galleriesCard.Card.RegularCard.VideoPreviewId != nil {
			var previewMedia media_entity.Media
			if err := tx.First(&previewMedia, *galleriesCard.Card.RegularCard.VideoPreviewId).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			galleriesCard.Card.RegularCard.VideoPreview = &previewMedia
		}
		if galleriesCard.Card.RegularCard != nil {
			if err := tx.Omit(clause.Associations).Create(galleriesCard.Card.RegularCard).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			for _, tag := range galleriesCard.Card.RegularCard.RegularCardsTags {

				if err := tx.Omit(clause.Associations).Create(tag).Error; err != nil {
					tx.Rollback()
					return nil, err
				}

			}
		}

		if galleriesCard.Card.VideoCard != nil {
			if err := tx.Omit(clause.Associations).Create(galleriesCard.Card.VideoCard).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if galleriesCard.Card.HtmlCard != nil {
			if err := tx.Omit(clause.Associations).Omit(clause.Associations).Create(galleriesCard.Card.HtmlCard).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if galleriesCard.Card.PhotoCard != nil && galleriesCard.Card.PhotoCard.PictureId != nil {
			var picture media_entity.Media
			if err := tx.First(&picture, *galleriesCard.Card.PhotoCard.PictureId).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			galleriesCard.Card.PhotoCard.Picture = &picture
		}

		if galleriesCard.Card.PhotoCard != nil {
			if err := tx.Omit(clause.Associations).Create(galleriesCard.Card.PhotoCard).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if galleriesCard.Card.FormCard != nil {
			if err := tx.Omit(clause.Associations).Create(galleriesCard.Card.FormCard).Error; err != nil {
				tx.Rollback()
				return nil, err
			}
			for _, tag := range galleriesCard.Card.FormCard.FormCardsTags {

				if err := tx.Omit(clause.Associations).Create(tag).Error; err != nil {
					tx.Rollback()
					return nil, err
				}

			}
		}

		if err := tx.Omit(clause.Associations).Create(&galleriesCard.Card).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := tx.Omit(clause.Associations).Create(&galleriesCard).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &gallery.ID, nil
}

func (r *GalleryRepository) Update(ctx context.Context, gallery *entity.Gallery) error {
	tx := r.db.WithContext(ctx).Begin()

	if err := tx.Omit(clause.Associations).Model(gallery).Where("id = ?", gallery.ID).
		Updates(
			gallery,
			//map[string]interface{}{
			//	"name":         gallery.Name,
			//	"code":         gallery.Code,
			//	"is_published": gallery.IsPublished,
			//	"hidden":       gallery.Hidden,
			//},
		).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, galleriesCard := range gallery.GalleriesCards {

		if galleriesCard.Card.VideoCard != nil && galleriesCard.Card.VideoCard.VideoId != nil {
			var videoMedia media_entity.Media
			if err := tx.First(&videoMedia, *galleriesCard.Card.VideoCard.VideoId).Error; err != nil {
				tx.Rollback()
				return err
			}
			galleriesCard.Card.VideoCard.Video = &videoMedia
		}

		if galleriesCard.Card.RegularCard != nil && galleriesCard.Card.RegularCard.VideoId != nil {
			var regularVideoMedia media_entity.Media
			if err := tx.First(&regularVideoMedia, *galleriesCard.Card.RegularCard.VideoId).Error; err != nil {
				tx.Rollback()
				return err
			}
			galleriesCard.Card.RegularCard.Video = &regularVideoMedia
		}

		if galleriesCard.Card.RegularCard != nil && galleriesCard.Card.RegularCard.VideoPreviewId != nil {
			var previewMedia media_entity.Media
			if err := tx.First(&previewMedia, *galleriesCard.Card.RegularCard.VideoPreviewId).Error; err != nil {
				tx.Rollback()
				return err
			}
			galleriesCard.Card.RegularCard.VideoPreview = &previewMedia
		}
		if galleriesCard.Card.RegularCard != nil {
			if err := tx.Omit(clause.Associations).Create(galleriesCard.Card.RegularCard).Error; err != nil {
				tx.Rollback()
				return err
			}
			for _, tag := range galleriesCard.Card.RegularCard.RegularCardsTags {

				if err := tx.Omit(clause.Associations).Create(tag).Error; err != nil {
					tx.Rollback()
					return err
				}

			}
		}

		if galleriesCard.Card.VideoCard != nil {
			if err := tx.Omit(clause.Associations).Create(galleriesCard.Card.VideoCard).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		if galleriesCard.Card.HtmlCard != nil {
			if err := tx.Omit(clause.Associations).Omit(clause.Associations).Create(galleriesCard.Card.HtmlCard).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		if galleriesCard.Card.PhotoCard != nil && galleriesCard.Card.PhotoCard.PictureId != nil {
			var picture media_entity.Media
			if err := tx.First(&picture, *galleriesCard.Card.PhotoCard.PictureId).Error; err != nil {
				tx.Rollback()
				return err
			}
			galleriesCard.Card.PhotoCard.Picture = &picture
		}

		if galleriesCard.Card.PhotoCard != nil {
			if err := tx.Omit(clause.Associations).Create(galleriesCard.Card.PhotoCard).Error; err != nil {
				tx.Rollback()
				return err
			}
		}

		if galleriesCard.Card.FormCard != nil {
			if err := tx.Omit(clause.Associations).Create(galleriesCard.Card.FormCard).Error; err != nil {
				tx.Rollback()
				return err
			}
			for _, tag := range galleriesCard.Card.FormCard.FormCardsTags {

				if err := tx.Omit(clause.Associations).Create(tag).Error; err != nil {
					tx.Rollback()
					return err
				}

			}
		}

		if err := tx.Omit(clause.Associations).Create(&galleriesCard.Card).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Omit(clause.Associations).Create(&galleriesCard).Error; err != nil {
			tx.Rollback()
			return err
		}

	}
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *GalleryRepository) UpdatePublish(ctx context.Context, galleryID uuid.UUID, publish *bool) error {
	err := r.db.WithContext(ctx).Model(&entity.Gallery{}).Where("id = ?", galleryID).
		Updates(
			map[string]interface{}{
				"is_published": *publish,
			},
		).Error
	return err
}

func (r *GalleryRepository) UpdateHidden(ctx context.Context, galleryID uuid.UUID, hidden *bool) error {
	err := r.db.WithContext(ctx).Model(&entity.Gallery{}).Where("id = ?", galleryID).
		Updates(
			map[string]interface{}{
				"hidden": *hidden,
			},
		).Error
	return err
}

func (r *GalleryRepository) PatchName(ctx context.Context, gallery *entity.Gallery) error {
	err := r.db.WithContext(ctx).Model(gallery).Where("id = ?", gallery.ID).
		Updates(
			map[string]interface{}{
				"name": gallery.Name,
			},
		).Error
	return err
}

func (r *GalleryRepository) PatchCardPosition(ctx context.Context, dto *actions.PatchCardPosition) error {
	tx := r.db.WithContext(ctx).Begin()

	err := r.UpdateCardsPositionAfterPatch(ctx, dto, tx)
	if err != nil {
		return err
	}

	err = tx.WithContext(ctx).Table("galleries_cards").
		Where("gallery_id = ?", dto.GalleryID).
		Where("card_id = ?", dto.CardID).
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

func (r *GalleryRepository) UpdateCardsPositionAfterPatch(
	ctx context.Context, dto *actions.PatchCardPosition, tx *gorm.DB,
) error {
	if *dto.OldPosition > *dto.NewPosition {
		err := tx.WithContext(ctx).Table("galleries_cards").Where("gallery_id = ?", dto.GalleryID).Where(
			"position >= ? and position < ?", dto.NewPosition, dto.OldPosition,
		).Update("position", gorm.Expr("position + 1")).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if *dto.OldPosition < *dto.NewPosition {
		err := tx.WithContext(ctx).Table("galleries_cards").Where("gallery_id = ?", dto.GalleryID).Where(
			"position <= ? and position > ?", dto.NewPosition, dto.OldPosition,
		).Update("position", gorm.Expr("position - 1")).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return nil
}

func (r *GalleryRepository) CreateGalleriesCards(ctx context.Context, galleriesCards []entity.GalleriesCards) (
	error, bool,
) {
	err := r.db.WithContext(ctx).Omit(clause.Associations).Create(&galleriesCards).Error

	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == UniqueViolationErr {
			return errors.BadRequest.Wrapf(pgErr, "card already linked to the gallery, detail: %s", pgErr.Detail), true
		}
	}

	return err, false
}

func (r *GalleryRepository) DeleteGalleriesCards(ctx context.Context, galleryID uuid.UUID, cardIDs []uuid.UUID) error {
	var galleriesCards []entity.GalleriesCards

	for _, cardId := range cardIDs {
		ok := r.CheckLink(ctx, galleryID, cardId)
		if !ok {
			return errors.NotFound.Newf("card with id: %s is not linked to the gallery", cardId.String())
		}
	}

	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).
		Where("gallery_id = ?", galleryID).
		Where("card_id in ?", cardIDs).
		Delete(&galleriesCards).Error
	if err != nil {
		return err
	}

	err = r.UpdateCardsPositionAfterDelete(ctx, galleryID, galleriesCards, err)

	return err
}

func (r *GalleryRepository) UpdateCardsPositionAfterDelete(
	ctx context.Context, galleryID uuid.UUID, galleriesCards []entity.GalleriesCards, err error,
) error {
	sort.Slice(
		galleriesCards, func(i, j int) bool {
			return galleriesCards[i].Position > galleriesCards[j].Position
		},
	)

	for _, gallery := range galleriesCards {
		err = r.db.WithContext(ctx).Table("galleries_cards").Where("gallery_id = ?", galleryID).Where(
			"position > ?", gallery.Position,
		).Update("position", gorm.Expr("position - 1")).Error
	}
	return err
}

func (r *GalleryRepository) GetLastPosition(ctx context.Context, galleryID uuid.UUID) (int, error) {
	var result int
	var count int64
	r.db.WithContext(ctx).Table("galleries_cards").Where(
		"gallery_id = ?", galleryID,
	).Count(&count)
	if count == 0 {
		return 0, nil
	}
	err := r.db.WithContext(ctx).Table("galleries_cards").Where(
		"gallery_id = ?", galleryID,
	).Select("max(position)").Row().Scan(&result)

	return result, err
}

func (r *GalleryRepository) DeleteList(ctx context.Context, galleryIds []uuid.UUID) error {
	err := r.UpdatePositionBeforeDelete(ctx, galleryIds)

	err = r.db.WithContext(ctx).Where("id in ?", galleryIds).Delete(entity.Gallery{}).Error
	return err
}

func (r *GalleryRepository) UpdatePositionBeforeDelete(ctx context.Context, galleryIds []uuid.UUID) error {
	var pagesGalleries []entity.PagesGalleries

	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Where(
		"gallery_id in ?", galleryIds,
	).Delete(&pagesGalleries).Error

	sort.Slice(
		pagesGalleries, func(i, j int) bool {
			return pagesGalleries[i].Position > pagesGalleries[j].Position
		},
	)

	for _, galleries := range pagesGalleries {
		err = r.db.WithContext(ctx).Table("pages_galleries").Where("pages_id = ?", galleries.PagesID).Where(
			"position > ?", galleries.Position,
		).Update("position", gorm.Expr("position - 1")).Error
	}
	return err
}

func (r *GalleryRepository) CheckLink(ctx context.Context, galleryID uuid.UUID, cardID uuid.UUID) bool {
	galleriesCard := &entity.GalleriesCards{}
	err := r.db.WithContext(ctx).Omit(clause.Associations).Model(entity.GalleriesCards{}).
		Where("gallery_id = ?", galleryID).
		Where("card_id = ?", cardID).
		First(galleriesCard).Error
	if err != nil || galleriesCard == nil {
		return false
	}
	return true
}
