package repositories

import (
	"context"
	"fmt"
	"github.com/aeroideaservices/focus/services/db/db_types/json"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/aeroideaservices/focus/media/plugin/actions"
	"github.com/aeroideaservices/focus/media/plugin/entity"
	"github.com/aeroideaservices/focus/services/errors"
)

// mediaRepository репозиторий медиа
type mediaRepository struct {
	db *gorm.DB
}

// NewMediaRepository конструктор
func NewMediaRepository(db *gorm.DB) actions.MediaRepository {
	return &mediaRepository{db: db}
}

// Create создание медиа
func (r mediaRepository) Create(ctx context.Context, medias ...entity.Media) error {
	err := r.db.WithContext(ctx).Create(medias).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error creating medias")
	}

	return nil
}

// Has проверка существования медиа
func (r mediaRepository) Has(ctx context.Context, id uuid.UUID) bool {
	err := r.db.WithContext(ctx).
		Select("id").
		Where("id = ?", id).
		First(&entity.Media{}).Error

	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// HasByFilter проверка существования медиа по фильтру
func (r mediaRepository) HasByFilter(ctx context.Context, filter actions.MediaFilter) bool {
	db := r.db.WithContext(ctx).Select("id")
	db = r.filterMedia(db, filter)
	err := db.First(&entity.Media{}).Error

	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// HasByFilterWithId проверка существования медиа по фильтру и возвращает id
func (r mediaRepository) HasByFilterWithId(ctx context.Context, filter actions.MediaFilter) (bool, uuid.UUID) {
	media := &entity.Media{}
	db := r.db.WithContext(ctx).Select("id")
	db = r.filterMedia(db, filter)
	err := db.First(media).Error

	return !errors.Is(err, gorm.ErrRecordNotFound), media.Id
}

// Get получение медиа по id
func (r mediaRepository) Get(ctx context.Context, id uuid.UUID) (*entity.Media, error) {
	media := &entity.Media{}
	err := r.db.WithContext(ctx).Where("id", id).First(media).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, actions.ErrMediaNotFound
	}
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting media")
	}

	return media, nil
}

// Update обновление медиа
func (r mediaRepository) Update(ctx context.Context, medias ...*actions.UpdateMediaDto) error {
	mediasTemplate := make([]string, len(medias))
	mediasValues := make([]any, 0, 7*len(medias))
	for i, media := range medias {
		mediasTemplate[i] = fmt.Sprintf("(?::uuid, ?, ?, ?, ?::uuid)")
		mediasValues = append(mediasValues, media.Id, media.Name, media.Filepath, media.Filename, media.FolderId)
	}

	err := r.db.WithContext(ctx).
		Exec(
			"WITH values (id, name, filepath, filename, folder_id)"+
				" AS (VALUES "+strings.Join(mediasTemplate, ", ")+")"+
				" UPDATE media"+
				" SET (name, filepath, filename, folder_id) = (v.name, v.filepath, v.filename, v.folder_id)"+
				" FROM values as v"+
				" WHERE v.id = media.id", mediasValues...,
		).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error updating medias")
	}

	return nil
}

// Delete удаление медиа
func (r mediaRepository) Delete(ctx context.Context, ids ...uuid.UUID) error {
	err := r.db.WithContext(ctx).
		Where("id IN (?)", ids).
		Delete(&entity.Media{}).
		Error
	if err != nil {
		return errors.NoType.Wrap(err, "error deleting medias")
	}

	return nil
}

func (r mediaRepository) filterMedia(db *gorm.DB, filter actions.MediaFilter) *gorm.DB {
	if filter.WithFolderId {
		db = db.Where("folder_id", filter.FolderId)
	}
	if filter.Name != "" {
		db = db.Where("name = ?", filter.Name)
	}
	if filter.Filename != "" {
		db = db.Where("filename = ?", filter.Filename)
	}
	if len(filter.Filenames) != 0 {
		db = db.Where("filename IN (?)", filter.Filenames)
	}
	if filter.Ext != "" {
		db = db.Where("ext = ?", filter.Ext)
	}

	return db
}

// GetShortList получение списка с краткой информацией по медиа
func (r mediaRepository) GetShortList(ctx context.Context, ids []uuid.UUID) ([]entity.Media, error) {
	var entities []entity.Media
	err := r.db.WithContext(ctx).
		Select("id, filepath, alt, title").
		Where("id IN (?)", ids).
		Find(&entities).
		Error
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting medias shorts")
	}

	return entities, nil
}

type gormScope func(*gorm.DB) *gorm.DB

// Count подсчет медиа
func (r mediaRepository) Count(ctx context.Context, filter actions.MediaFilter) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entity.Media{}).
		Scopes(getMediaFilterScopes(filter)).
		Select("id").Distinct().
		Count(&count).
		Error
	if err != nil {
		return 0, errors.NoType.Wrap(err, "error getting medias count")
	}

	return int(count), nil
}

func (r mediaRepository) UpdateSubtitles(ctx context.Context, mediaId uuid.UUID, subtitle json.JSONB) error {
	err := r.db.WithContext(ctx).
		Model(&entity.Media{}).
		Where("id = ?", mediaId).
		Update("subtitle", subtitle).
		Error
	if err != nil {
		return errors.NoType.Wrap(err, "error updating media subtitle")
	}
	return nil
}

func getMediaFilterScopes(filter actions.MediaFilter) gormScope {
	return func(db *gorm.DB) *gorm.DB {
		if filter.WithFolderId {
			db = db.Where("folder_id", filter.FolderId)
		}
		if filter.Name != "" {
			db = db.Where("name = ?", filter.Name)
		}
		if filter.Filename != "" {
			db = db.Where("filename = ?", filter.Filename)
		}
		if len(filter.Filenames) != 0 {
			db = db.Where("filename IN (?)", filter.Filenames)
		}
		if filter.Ext != "" {
			db = db.Where("ext = ?", filter.Ext)
		}
		if len(filter.InIds) != 0 {
			inIds := make([]any, len(filter.InIds))
			for i := range filter.InIds {
				inIds[i] = filter.InIds[i]
			}
			db = db.Where(clause.IN{Column: clause.Column{Table: clause.CurrentTable, Name: "id"}, Values: inIds})
		}

		return db
	}
}
