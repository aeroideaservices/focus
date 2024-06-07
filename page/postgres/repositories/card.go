package repositories

import (
	"context"
	media_entity "github.com/aeroideaservices/focus/media/plugin/entity"
	"github.com/aeroideaservices/focus/page/plugin/entity"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{
		db: db,
	}
}

const (
	UniqueViolationErr = "23505"
)

func (r *CardRepository) GetListWithSearch(ctx context.Context, searchValue string, name string) (
	[]entity.Card, error,
) {
	var cards []entity.Card

	db := r.db.WithContext(ctx).Model(entity.Card{}).
		Preload("HtmlCard").
		Preload("VideoCard").Preload("VideoCard.Video").Preload("VideoCard.VideoLite").Preload("VideoCard.VideoPreview").Preload("VideoCard.VideoPreviewBlur").
		//Preload(
		//	"RegularCard", func(db *gorm.DB) *gorm.DB {
		//		return db.Where("lower(preview_text) LIKE lower(?)", "%"+searchValue+"%")
		//	},
		//).
		Preload("RegularCard.Video").Preload("RegularCard.VideoLite").Preload("RegularCard.VideoPreview").Preload("RegularCard.VideoPreviewBlur").Preload("RegularCard.User").
		Preload("RegularCard.RegularCardsTags.Tag").
		Preload("RegularCard.User.Picture").
		Preload("PhotoCard").Preload("PhotoCard.Picture").
		Preload("FormCard").Preload("FormCard.Background").Preload("FormCard.User").Preload("FormCard.User.Picture").
		Preload("FormCard.FormCardsTags.Tag").Preload("FormCard.Form").
		Where("lower(type) LIKE lower(?)", "%"+searchValue+"%")
	if name != "" {
		db.Where("name ilike ?", "%"+name+"%")
	}
	//Or("lower(RegularCard.preview_text) LIKE lower(?)", "%"+searchValue+"%").
	//Or("lower(HtmlCard.html) LIKE lower(?)", "%"+searchValue+"%")
	err := db.Find(&cards).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrapf(err, "cards with searchValue %s not found", searchValue)
	}
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *CardRepository) GetById(ctx context.Context, cardId uuid.UUID) (*entity.Card, error) {
	card := &entity.Card{}

	db := r.db.WithContext(ctx).Model(entity.Card{}).Where("id", cardId).Preload("HtmlCard").
		Preload("VideoCard").Preload("VideoCard.Video").Preload("VideoCard.VideoLite").Preload("VideoCard.VideoPreview").Preload("VideoCard.VideoPreviewBlur").
		Preload("RegularCard").Preload("RegularCard.Video").Preload("RegularCard.VideoLite").Preload("RegularCard.VideoPreview").Preload("RegularCard.VideoPreviewBlur").
		Preload("RegularCard.User").Preload("RegularCard.User.Picture").
		Preload("RegularCard.RegularCardsTags.Tag").
		Preload("PhotoCard").Preload("PhotoCard.Picture").
		Preload("FormCard").Preload("FormCard.Background").Preload("FormCard.User").Preload("FormCard.User.Picture").
		Preload("FormCard.FormCardsTags.Tag").Preload("FormCard.Form")

	err := db.First(card).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrapf(err, "card with id %s not found", cardId)
	}
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (r *CardRepository) GetByIdWithoutAssociate(ctx context.Context, cardId uuid.UUID) (*entity.Card, error) {
	card := &entity.Card{}

	db := r.db.WithContext(ctx).Omit(clause.Associations).Model(entity.Card{}).Where("id", cardId)
	err := db.First(card).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrapf(err, "card with id %s not found", cardId)
	}
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (r *CardRepository) GetListByGalleryId(ctx context.Context, galleryId uuid.UUID) ([]entity.Card, error) {
	var cards []entity.Card
	db := r.db.WithContext(ctx).Select("cards.*, galleries_cards.position").Table("galleries_cards").
		Joins("left join cards on galleries_cards.card_id=cards.id").
		Where("gallery_id=?", galleryId).
		Preload("HtmlCard").
		Preload("VideoCard").Preload("VideoCard.Video").Preload("VideoCard.VideoLite").Preload("VideoCard.VideoPreview").Preload("VideoCard.VideoPreviewBlur").
		Preload("RegularCard").Preload("RegularCard.Video").Preload("RegularCard.VideoLite").Preload("RegularCard.VideoPreview").Preload("RegularCard.VideoPreviewBlur").Preload("RegularCard.User").
		Preload("RegularCard.RegularCardsTags.Tag").
		Preload("PhotoCard").Preload("PhotoCard.Picture").
		Preload("FormCard").Preload("FormCard.Background").Preload("FormCard.User").Preload("FormCard.User.Picture").
		Preload("FormCard.FormCardsTags.Tag").Preload("FormCard.Form").
		Find(&cards)

	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, db.Error
		}
	}

	return cards, nil
}

//func (r *CardRepository) Create(
//	ctx context.Context, card *entity.Card, galleriesCards []entity.GalleriesCards,
//) (*uuid.UUID, error) {
//	tx := r.db.WithContext(ctx).Begin()
//
//	if card.VideoCard != nil && card.VideoCard.VideoId != nil {
//		var videoMedia media_entity.Media
//		if err := tx.First(&videoMedia, *card.VideoCard.VideoId).Error; err != nil {
//			tx.Rollback()
//			return nil, err
//		}
//		card.VideoCard.Video = &videoMedia
//	}
//
//	if card.RegularCard != nil && card.RegularCard.VideoId != nil {
//		var regularVideoMedia media_entity.Media
//		if err := tx.First(&regularVideoMedia, *card.RegularCard.VideoId).Error; err != nil {
//			tx.Rollback()
//			return nil, err
//		}
//		card.RegularCard.Video = &regularVideoMedia
//	}
//
//	if card.RegularCard != nil && card.RegularCard.VideoPreviewId != nil {
//		var previewMedia media_entity.Media
//		if err := tx.First(&previewMedia, *card.RegularCard.VideoPreviewId).Error; err != nil {
//			tx.Rollback()
//			return nil, err
//		}
//		card.RegularCard.VideoPreview = &previewMedia
//	}
//	if card.RegularCard != nil {
//		if err := tx.Omit(clause.Associations).Create(card.RegularCard).Error; err != nil {
//			tx.Rollback()
//			return nil, err
//		}
//		for _, tag := range card.RegularCard.RegularCardsTags {
//
//			if err := tx.Omit(clause.Associations).Create(tag).Error; err != nil {
//				tx.Rollback()
//				return nil, err
//			}
//
//		}
//	}
//
//	if card.VideoCard != nil {
//		if err := tx.Omit(clause.Associations).Create(card.VideoCard).Error; err != nil {
//			tx.Rollback()
//			return nil, err
//		}
//	}
//
//	if card.HtmlCard != nil {
//		if err := tx.Omit(clause.Associations).Omit(clause.Associations).Create(card.HtmlCard).Error; err != nil {
//			tx.Rollback()
//			return nil, err
//		}
//	}
//
//	if card.PhotoCard != nil {
//		if err := tx.Omit(clause.Associations).Omit(clause.Associations).Create(card.PhotoCard).Error; err != nil {
//			tx.Rollback()
//			return nil, err
//		}
//	}
//
//	if card.FormCard != nil {
//		if err := tx.Omit(clause.Associations).Omit(clause.Associations).Create(card.FormCard).Error; err != nil {
//			tx.Rollback()
//			return nil, err
//		}
//		for _, tag := range card.FormCard.FormCardsTags {
//
//			if err := tx.Omit(clause.Associations).Create(tag).Error; err != nil {
//				tx.Rollback()
//				return nil, err
//			}
//
//		}
//	}
//
//	if err := tx.Omit(clause.Associations).Create(&card).Error; err != nil {
//		tx.Rollback()
//		return nil, err
//	}
//
//	if err := tx.Commit().Error; err != nil {
//		return nil, err
//	}
//
//	if len(galleriesCards) != 0 {
//		err := r.db.WithContext(ctx).Omit(clause.Associations).Create(&galleriesCards).Error
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return &card.ID, nil
//}

func (r *CardRepository) Create(
	ctx context.Context, card *entity.Card, galleriesCards []entity.GalleriesCards,
) (*uuid.UUID, error) {
	tx := r.db.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var err error

	if card.VideoCard != nil && card.VideoCard.VideoId != nil {
		if videoMedia, err := r.getMedia(tx, *card.VideoCard.VideoId); err != nil {
			return nil, err
		} else {
			card.VideoCard.Video = videoMedia
		}
	}

	if card.RegularCard != nil {
		if err = r.createRegularCard(tx, card); err != nil {
			return nil, err
		}
	}

	if card.VideoCard != nil {
		if err = r.createVideoCard(tx, card); err != nil {
			return nil, err
		}
	}

	if card.HtmlCard != nil {
		if err = r.createHtmlCard(tx, card); err != nil {
			return nil, err
		}
	}

	if card.PhotoCard != nil {
		if err = r.createPhotoCard(tx, card); err != nil {
			return nil, err
		}
	}

	if card.FormCard != nil {
		if err = r.createFormCard(tx, card); err != nil {
			return nil, err
		}
	}

	if err = r.createCard(tx, card); err != nil {
		return nil, err
	}

	if err = tx.Commit().Error; err != nil {
		return nil, err
	}

	if len(galleriesCards) != 0 {
		if err = r.db.WithContext(ctx).Omit(clause.Associations).Create(&galleriesCards).Error; err != nil {
			return nil, err
		}
	}

	return &card.ID, nil
}

func (r *CardRepository) getMedia(tx *gorm.DB, mediaId uuid.UUID) (*media_entity.Media, error) {
	var media media_entity.Media
	if err := tx.First(&media, mediaId).Error; err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *CardRepository) createRegularCard(tx *gorm.DB, card *entity.Card) error {
	if card.RegularCard.VideoId != nil {
		videoMedia, err := r.getMedia(tx, *card.RegularCard.VideoId)
		if err != nil {
			return err
		}
		card.RegularCard.Video = videoMedia
	}

	if card.RegularCard.VideoPreviewId != nil {
		previewMedia, err := r.getMedia(tx, *card.RegularCard.VideoPreviewId)
		if err != nil {
			return err
		}
		card.RegularCard.VideoPreview = previewMedia
	}

	if err := tx.Omit(clause.Associations).Create(card.RegularCard).Error; err != nil {
		return err
	}

	for _, tag := range card.RegularCard.RegularCardsTags {
		if err := tx.Omit(clause.Associations).Create(tag).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *CardRepository) createVideoCard(tx *gorm.DB, card *entity.Card) error {
	if err := tx.Omit(clause.Associations).Create(card.VideoCard).Error; err != nil {
		return err
	}
	return nil
}

func (r *CardRepository) createHtmlCard(tx *gorm.DB, card *entity.Card) error {
	if err := tx.Omit(clause.Associations).Create(card.HtmlCard).Error; err != nil {
		return err
	}
	return nil
}

func (r *CardRepository) createPhotoCard(tx *gorm.DB, card *entity.Card) error {
	if err := tx.Omit(clause.Associations).Create(card.PhotoCard).Error; err != nil {
		return err
	}
	return nil
}

func (r *CardRepository) createFormCard(tx *gorm.DB, card *entity.Card) error {
	if err := tx.Omit(clause.Associations).Create(card.FormCard).Error; err != nil {
		return err
	}
	return nil
}

func (r *CardRepository) createCard(tx *gorm.DB, card *entity.Card) error {
	if err := tx.Omit(clause.Associations).Create(card).Error; err != nil {
		return err
	}
	return nil
}

func (r *CardRepository) Update(
	ctx context.Context, card *entity.Card, galleriesCards []entity.GalleriesCards,
) error {
	internalCardId, err := r.getInternalId(ctx, card)
	if err != nil {
		return err
	}

	switch card.Type {
	default:
		err := r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).
			Updates(
				card,
			).Error
		if err != nil {
			return err
		}

	case "html":
		card.HtmlCard.ID = uuid.Must(uuid.Parse(internalCardId))
		card.HtmlCardId = &card.HtmlCard.ID

		err = r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).
			Updates(
				card,
			).Error
		if err != nil {
			return err
		}
	case "video":
		card.VideoCard.ID = uuid.Must(uuid.Parse(internalCardId))
		card.VideoCardId = &card.VideoCard.ID

		err = r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).
			Updates(
				card,
			).Error
		if err != nil {
			return err
		}
	case "photo":
		card.PhotoCard.ID = uuid.Must(uuid.Parse(internalCardId))
		card.PhotoCardId = &card.PhotoCard.ID

		err = r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).
			Updates(
				card,
			).Error
		if err != nil {
			return err
		}
	case "form":
		card.FormCard.ID = uuid.Must(uuid.Parse(internalCardId))
		card.FormCardId = &card.FormCard.ID

		for i, _ := range card.FormCard.FormCardsTags {
			card.FormCard.FormCardsTags[i].FormCardID = card.FormCard.ID
		}

		err = r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Omit("FormCard.FormCardsTags").
			Updates(
				card,
				//map[string]interface{}{
				//	"name":         gallery.Name,
				//	"code":         gallery.Code,
				//	"is_published": gallery.IsPublished,
				//	"hidden":       gallery.Hidden,
				//},
			).Error

		if err != nil {
			return err
		}
		err = r.db.WithContext(ctx).Omit(clause.Associations).Create(card.FormCard.FormCardsTags).Error

	case "regular":
		card.RegularCard.ID = uuid.Must(uuid.Parse(internalCardId))
		card.RegularCardId = &card.RegularCard.ID

		for i, _ := range card.RegularCard.RegularCardsTags {
			card.RegularCard.RegularCardsTags[i].RegularCardID = card.RegularCard.ID
		}

		err = r.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Omit("RegularCard.RegularCardsTags").
			Updates(
				card,
				//map[string]interface{}{
				//	"name":         gallery.Name,
				//	"code":         gallery.Code,
				//	"is_published": gallery.IsPublished,
				//	"hidden":       gallery.Hidden,
				//},
			).Error

		if err != nil {
			return err
		}
		err = r.db.WithContext(ctx).Omit(clause.Associations).Create(card.RegularCard.RegularCardsTags).Error
	}

	if len(galleriesCards) != 0 {
		err := r.db.WithContext(ctx).Omit(clause.Associations).Create(&galleriesCards).Error
		if err != nil {
			return err
		}
	}

	//err := r.db.WithContext(ctx).Omit(clause.Associations).Create(&galleriesCards).Error

	//err = r.db.WithContext(ctx).Omit(clause.Associations).Create(card.RegularCard.RegularCardsTags).Error

	return nil
}

func (r *CardRepository) getInternalId(ctx context.Context, card *entity.Card) (string, error) {
	var htmlCardId string
	selectString := card.Type + "_card_id"
	err := r.db.WithContext(ctx).Select(selectString).Model(entity.Card{}).Where("id = ?", card.ID).
		Find(&htmlCardId).Error
	if err != nil {
		return "", err
	}
	return htmlCardId, err
}

func (r *CardRepository) Delete(ctx context.Context, cardId uuid.UUID) error {

	err := r.UpdatePositionBeforeDelete(ctx, cardId)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Where("id = ?", cardId).Delete(entity.Card{}).Error

	return err
}

func (r *CardRepository) UpdatePositionBeforeDelete(ctx context.Context, cardId uuid.UUID) error {
	var galleriesCard []entity.GalleriesCards

	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Where("card_id = ?", cardId).Delete(&galleriesCard).Error

	for _, cards := range galleriesCard {
		err = r.db.WithContext(ctx).Table("galleries_cards").Where("gallery_id = ?", cards.GalleryID).Where(
			"position > ?", cards.Position,
		).Update("position", gorm.Expr("position - 1")).Error
	}
	return err
}

func (r *CardRepository) PatchUser(ctx context.Context, card *entity.Card) error {
	//var cardType string
	newUserId := card.RegularCard.UserId
	err := r.db.WithContext(ctx).Select("type", "regular_card_id", "form_card_id").Model(entity.Card{}).Where(
		"id = ?", card.ID,
	).Find(card).Error

	if err != nil {
		return err
	}

	if card.RegularCardId != nil {
		err = r.db.WithContext(ctx).Model(entity.RegularCard{}).Where("id = ?", card.RegularCardId).
			Updates(
				map[string]interface{}{
					"user_id": newUserId,
				},
			).Error

	}
	if card.FormCardId != nil {
		err = r.db.WithContext(ctx).Model(entity.FormCard{}).Where("id = ?", card.FormCardId).
			Updates(
				map[string]interface{}{
					"user_id": newUserId,
				},
			).Error

	}

	return err
}

func (r *CardRepository) PatchPreviewText(ctx context.Context, card *entity.Card) error {
	//var cardType string

	err := r.db.WithContext(ctx).Select("type", "regular_card_id").Model(entity.Card{}).Where(
		"id = ?", card.ID,
	).Find(card).Error

	if err != nil {
		return err
	}

	if card.RegularCardId != nil {
		err = r.db.WithContext(ctx).Model(entity.RegularCard{}).Where("id = ?", card.RegularCardId).
			Updates(
				map[string]interface{}{
					"preview_text": card.RegularCard.PreviewText,
				},
			).Error

	}
	return err
}

func (r *CardRepository) PatchDetailText(ctx context.Context, card *entity.Card) error {
	//var cardType string
	//newUserId := card.RegularCard.UserId
	err := r.db.WithContext(ctx).Select("type", "regular_card_id").Model(entity.Card{}).Where(
		"id = ?", card.ID,
	).Find(card).Error

	if err != nil {
		return err
	}

	if card.RegularCardId != nil {
		err = r.db.WithContext(ctx).Model(entity.RegularCard{}).Where("id = ?", card.RegularCardId).
			Updates(
				map[string]interface{}{
					"detail_text": card.RegularCard.DetailText,
				},
			).Error

	}
	return err
}

func (r *CardRepository) PatchLearnMoreUrl(ctx context.Context, card *entity.Card) error {
	err := r.db.WithContext(ctx).Select("type", "regular_card_id", "form_card_id").Model(entity.Card{}).Where(
		"id = ?", card.ID,
	).Find(card).Error

	if err != nil {
		return err
	}

	if card.RegularCardId != nil {
		return r.db.WithContext(ctx).Model(entity.RegularCard{}).Where("id = ?", card.RegularCardId).
			Updates(
				map[string]interface{}{
					"learn_more_url": card.RegularCard.LearnMoreUrl,
				},
			).Error
	}
	if card.FormCardId != nil {
		return r.db.WithContext(ctx).Model(entity.FormCard{}).Where("id = ?", card.FormCardId).
			Updates(
				map[string]interface{}{
					"learn_more_url": card.FormCard.LearnMoreUrl,
				},
			).Error
	}
	return err
}

func (r *CardRepository) PatchTags(ctx context.Context, cardId uuid.UUID, tags []entity.Tag) error {
	//regularCardId, err := r.GetRegularCardId(ctx, cardId)
	//if err != nil {
	//	return err
	//}
	var err error
	for _, tag := range tags {
		err = r.db.WithContext(ctx).Model(entity.Tag{}).Where("id = ?", tag.ID).Updates(&tag).Error
		if err != nil {
			return err
		}
	}
	//err = r.db.WithContext(ctx).Model(entity.Tag{}).Updates(&tags).Error

	return err
}

func (r *CardRepository) CreateRegularTags(ctx context.Context, cardId uuid.UUID, tagIds []uuid.UUID) (error, bool) {
	card, err := r.GetInternalCardId(ctx, cardId)
	if err != nil {
		return err, false
	}

	if card.RegularCardId != nil {
		var regularCardsTags []entity.RegularCardsTags
		for _, id := range tagIds {
			regularCardsTags = append(
				regularCardsTags, entity.RegularCardsTags{
					RegularCardID: *card.RegularCardId,
					TagID:         id,
				},
			)
		}
		err = r.db.WithContext(ctx).Omit(clause.Associations).Create(&regularCardsTags).Error
	}
	if card.FormCardId != nil {
		var formCardsTags []entity.FormCardsTags
		for _, id := range tagIds {
			formCardsTags = append(
				formCardsTags, entity.FormCardsTags{
					FormCardID: *card.FormCardId,
					TagID:      id,
				},
			)
		}
		err = r.db.WithContext(ctx).Omit(clause.Associations).Create(&formCardsTags).Error
	}

	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == UniqueViolationErr {
			return errors.BadRequest.Wrapf(pgErr, "tag already linked to the card, detail: %s", pgErr.Detail), true
		}
	}

	return err, false
}

func (r *CardRepository) DeleteRegularTags(ctx context.Context, cardId uuid.UUID, tagIds []uuid.UUID) error {
	card, err := r.GetInternalCardId(ctx, cardId)
	if err != nil {
		return err
	}

	for _, tagId := range tagIds {
		ok := r.CheckLink(ctx, card, tagId)
		if !ok {
			return errors.NotFound.Newf("tag with id: %s is not linked to the card", tagId.String())
		}
	}

	if card.RegularCardId != nil {
		err = r.db.WithContext(ctx).Clauses(clause.Returning{}).
			Where("regular_card_id = ?", card.RegularCardId).
			Where("tag_id in ?", tagIds).
			Delete(&entity.RegularCardsTags{}).Error
	}
	if card.FormCardId != nil {
		err = r.db.WithContext(ctx).Clauses(clause.Returning{}).
			Where("form_card_id = ?", card.FormCardId).
			Where("tag_id in ?", tagIds).
			Delete(&entity.FormCardsTags{}).Error
	}

	return err
}

func (r *CardRepository) GetRegularCardId(ctx context.Context, cardId uuid.UUID) (string, error) {
	var regularCardId string
	err := r.db.WithContext(ctx).Select("regular_card_id").Model(&entity.Card{}).Where("id = ?", cardId).
		Find(&regularCardId).Error
	return regularCardId, err
}

func (r *CardRepository) GetInternalCardId(ctx context.Context, cardId uuid.UUID) (*entity.Card, error) {
	card := &entity.Card{}
	err := r.db.WithContext(ctx).Select("type", "regular_card_id", "form_card_id").Model(entity.Card{}).Where(
		"id = ?", cardId,
	).Find(card).Error
	return card, err
}

func (r *CardRepository) UpdateInverted(ctx context.Context, cardID uuid.UUID, inverted *bool) error {
	regularCardId, err := r.GetRegularCardId(ctx, cardID)
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Model(&entity.RegularCard{}).Where("id = ?", regularCardId).
		Updates(
			map[string]interface{}{
				"inverted": *inverted,
			},
		).Error
	return err
}

func (r *CardRepository) UpdatePublish(ctx context.Context, cardID uuid.UUID, publish *bool) error {
	err := r.db.WithContext(ctx).Model(&entity.Card{}).Where("id = ?", cardID).
		Updates(
			map[string]interface{}{
				"is_published": *publish,
			},
		).Error
	return err
}

func (r *CardRepository) GetLastPositionInGalley(ctx context.Context, galleryID uuid.UUID) (int, error) {
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

func (r *CardRepository) CheckLink(ctx context.Context, card *entity.Card, tagID uuid.UUID) bool {

	if card.RegularCardId != nil {
		regularTags := &entity.RegularCardsTags{}
		err := r.db.WithContext(ctx).Omit(clause.Associations).Model(entity.RegularCardsTags{}).
			Where("regular_card_id = ?", card.RegularCardId).
			Where("tag_id = ?", tagID).
			First(regularTags).Error
		if err != nil || regularTags == nil {
			return false
		}
	}
	if card.FormCardId != nil {
		formTags := &entity.FormCardsTags{}
		err := r.db.WithContext(ctx).Omit(clause.Associations).Model(entity.FormCardsTags{}).
			Where("form_card_id = ?", card.FormCardId).
			Where("tag_id = ?", tagID).
			First(formTags).Error
		if err != nil || formTags == nil {
			return false
		}
	}
	return true
}
