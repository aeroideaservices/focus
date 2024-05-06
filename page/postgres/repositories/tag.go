package repositories

import (
	"context"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"pages/pkg/page/plugin/entity"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{
		db: db,
	}
}

func (r *TagRepository) UpdateIsDetailLink(ctx context.Context, tagID uuid.UUID, value *bool) error {
	err := r.db.WithContext(ctx).Model(&entity.Tag{}).Where("id = ?", tagID).
		Updates(
			map[string]interface{}{
				"is_primary": *value,
			},
		).Error
	return err
}

func (r *TagRepository) GetListWithSearch(ctx context.Context, searchValue string) (
	[]entity.Tag, map[uuid.UUID][]uuid.UUID, error,
) {
	var tags []entity.Tag

	db := r.db.WithContext(ctx).Model(entity.Tag{}).
		Where("lower(text) LIKE lower(?)", "%"+searchValue+"%")
	err := db.Find(&tags).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, errors.NotFound.Wrapf(err, "tags with searchValue %s not found", searchValue)
	}
	if err != nil {
		return nil, nil, err
	}

	tagCardsIds := make(map[uuid.UUID][]uuid.UUID, len(tags))

	for _, tag := range tags {
		cardIds, err := r.GetCardIdsByTagId(ctx, tag.ID)
		if err != nil {
			return nil, nil, err
		}
		tagCardsIds[tag.ID] = cardIds
	}

	return tags, tagCardsIds, nil
}

func (r *TagRepository) GetById(ctx context.Context, tagId uuid.UUID) (*entity.Tag, error) {
	tag := &entity.Tag{}

	db := r.db.WithContext(ctx).Model(entity.Tag{}).Where("id", tagId)

	err := db.First(tag).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.NotFound.Wrapf(err, "tag with id %s not found", tagId)
	}
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *TagRepository) GetCardIdsByTagId(ctx context.Context, tagId uuid.UUID) ([]uuid.UUID, error) {
	var regularCardIds []uuid.UUID
	var formCardIds []uuid.UUID
	var cardIds []uuid.UUID

	err := r.db.WithContext(ctx).Select("regular_card_id").Model(entity.RegularCardsTags{}).Where(
		"tag_id = ?", tagId,
	).Scan(&regularCardIds).Error

	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Select("form_card_id").Model(entity.FormCardsTags{}).Where(
		"tag_id = ?", tagId,
	).Scan(&formCardIds).Error

	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Select("id").Model(entity.Card{}).Where(
		"regular_card_id in ?", regularCardIds,
	).Or("form_card_id in ?", formCardIds).Scan(&cardIds).Error

	if err != nil {
		return nil, err
	}
	return cardIds, nil
}

func (r *TagRepository) Create(ctx context.Context, tag *entity.Tag) (*uuid.UUID, error) {
	err := r.db.WithContext(ctx).Create(tag).Error
	return &tag.ID, err
}

func (r *TagRepository) Update(ctx context.Context, tag *entity.Tag) error {
	err := r.db.WithContext(ctx).Updates(tag).Error
	return err
}

func (r *TagRepository) Delete(ctx context.Context, tagId uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("id = ?", tagId).Delete(entity.Tag{}).Error
	return err
}
