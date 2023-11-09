package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/WinterYukky/gorm-extra-clause-plugin/exclause"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/aeroideaservices/focus/media/plugin/actions"
	"github.com/aeroideaservices/focus/media/plugin/entity"
	"github.com/aeroideaservices/focus/services/errors"
)

// folderRepository репозиторий папок
type folderRepository struct {
	db *gorm.DB
}

// NewFolderRepository конструктор
func NewFolderRepository(db *gorm.DB) actions.FolderRepository {
	return &folderRepository{db: db}
}

// Has проверка существования папки по id
func (r folderRepository) Has(ctx context.Context, id uuid.UUID) bool {
	err := r.db.WithContext(ctx).
		Select("id").
		Where("id = ?", id).
		First(&entity.Folder{}).Error

	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// HasByFilter проверка существования папки по фильтру
func (r folderRepository) HasByFilter(ctx context.Context, filter actions.Filter) bool {
	db := r.db.WithContext(ctx).Select("id")
	db = r.filterFolder(db, filter)
	err := db.First(&entity.Folder{}).Error

	return !errors.Is(err, gorm.ErrRecordNotFound)
}

// Get получение папки по id
func (r folderRepository) Get(ctx context.Context, id uuid.UUID) (*entity.Folder, error) {
	folder := &entity.Folder{}
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(folder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, actions.ErrFolderNotFound
	}
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting folder")
	}

	return folder, nil
}

// GetWithSize получение папки с размером содержимого
func (r folderRepository) GetWithSize(ctx context.Context, id uuid.UUID) (*actions.FolderDetail, error) {
	res := &actions.FolderDetail{}

	err := r.db.Table(`
		(WITH RECURSIVE tree(id, folder_id, name, size) AS (
			SELECT id, folder_id, name, 
				   (SELECT coalesce(SUM(media.size), 0)::bigint FROM media WHERE media.folder_id = folders.id) as size
			FROM folders
			UNION ALL
			SELECT fol.id, fol.folder_id, fol.name, t.size
			FROM folders fol
			inner join tree t On fol.id = t.folder_id
		)
		SELECT id, name, folder_id, SUM(size) as size
		FROM tree
		GROUP BY id, name, folder_id) AS f
		`).
		WithContext(ctx).
		Where("f.id", id).
		First(&res).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, actions.ErrFolderNotFound
	}
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting folder with size")
	}

	return res, nil
}

// Create создание папки
func (r folderRepository) Create(ctx context.Context, folders ...*entity.Folder) error {
	err := r.db.WithContext(ctx).Create(folders).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error creating folder")
	}

	return nil
}

// Update обновление папки
func (r folderRepository) Update(ctx context.Context, folder *entity.Folder) error {
	err := r.db.Model(folder).WithContext(ctx).
		Updates(map[string]any{
			"name":      folder.Name,
			"folder_id": folder.FolderId,
		}).
		Error
	if err != nil {
		return errors.NoType.Wrap(err, "error updating folder")
	}

	return nil
}

// Delete удаление папки
func (r folderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Debug().Where("id = ?", id).Delete(&entity.Folder{}).Error
	if err != nil {
		return errors.NoType.Wrap(err, "error deleting folder")
	}

	return nil
}

// List получение списка папок
func (r folderRepository) List(ctx context.Context, filter actions.Filter) ([]*entity.Folder, error) {
	var res []*entity.Folder
	err := r.db.WithContext(ctx).Debug().
		Scopes(r.filterScopes(filter)).
		Find(&res).Error
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error listing folder")
	}

	return res, nil
}

// GetFolderPath получение пути папки
func (r folderRepository) GetFolderPath(ctx context.Context, id uuid.UUID) (folderPath string, err error) {
	var folderNames []string
	err = r.db.WithContext(ctx).Raw(
		`WITH RECURSIVE parent_folders (id, folder_id, name, folder_level) AS (
				SELECT id, folder_id, name, 1 folder_level
				FROM folders
				WHERE id = ?
				UNION ALL
				SELECT f.id, f.folder_id, f.name, folder_level + 1
				FROM folders f INNER JOIN parent_folders pf
			 	ON f.id = pf.folder_id
			)
			SELECT name FROM parent_folders
			ORDER BY folder_level DESC`, id,
	).Scan(&folderNames).Error
	if err != nil {
		return "", errors.NoType.Wrap(err, "error getting folder path")
	}

	folderPath = strings.Join(folderNames, "/")

	return folderPath, nil
}

// GetFolderMediaFilePaths получение путей медиа файлов
func (r folderRepository) GetFolderMediaFilePaths(ctx context.Context, id *uuid.UUID) (mediaFilepath []string, err error) {
	err = r.db.WithContext(ctx).Debug().Raw(
		`WITH RECURSIVE parent_folders (id, folder_id) AS (
				SELECT id, folder_id
				FROM folders
				WHERE id = ?
				UNION ALL
				SELECT f.id, f.folder_id
				FROM folders f
				INNER JOIN parent_folders pf
				ON f.folder_id = pf.id
			)
			SELECT DISTINCT filepath
			FROM media
			INNER JOIN parent_folders on parent_folders.id = media.folder_id`, id,
	).Scan(&mediaFilepath).Error
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting media filepath")
	}

	return mediaFilepath, nil
}

// GetAllFolderMedias получение всех медиа
func (r folderRepository) GetAllFolderMedias(ctx context.Context, id uuid.UUID) ([]*actions.UpdateMediaDto, error) {
	var medias []*actions.UpdateMediaDto
	err := r.db.WithContext(ctx).Debug().Raw(
		`WITH RECURSIVE parent_folders (id, folder_id, folder_path) AS (
				SELECT id, folder_id, name as folder_path
				FROM folders
				WHERE folder_id IS NULL
				UNION ALL
				SELECT f.id, f.folder_id, pf.folder_path || '/' || f.name
				FROM folders f
				INNER JOIN parent_folders pf
				ON f.folder_id = pf.id
			)
			SELECT DISTINCT media.id, media.name, media.filename, media.filepath, media.folder_id, parent_folders.folder_path || '/' || media.filename as new_filepath
			FROM media
			INNER JOIN parent_folders on parent_folders.id = media.folder_id
			WHERE media.folder_id = ?`, id,
	).Scan(&medias).Error
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting folder medias")
	}

	return medias, nil
}

// HasSubFolder проверка существования поддиректории
func (r folderRepository) HasSubFolder(ctx context.Context, id uuid.UUID, subFolderId *uuid.UUID) (bool, error) {
	var hasSubFolder bool
	err := r.db.WithContext(ctx).Raw(
		`WITH RECURSIVE parent_folders (id, folder_id) AS (
				SELECT id, folder_id
				FROM folders
				WHERE id = ?
				UNION ALL
				SELECT f.id, f.folder_id
				FROM folders f
				INNER JOIN parent_folders pf
				ON f.folder_id = pf.id
			)
			SELECT COUNT(1) > 0 count
			FROM parent_folders
			WHERE id = ?`, id, subFolderId,
	).Scan(&hasSubFolder).Error
	if err != nil {
		return false, errors.NoType.Wrap(err, "error checking if has sub folder")
	}

	return hasSubFolder, nil
}

// GetFoldersTree получение дерева директорий
func (r folderRepository) GetFoldersTree(ctx context.Context) ([]*actions.FolderResponse, error) {
	var entities []*actions.FolderResponse

	err := r.db.WithContext(ctx).
		Model(entity.Folder{}).
		Raw(`
		WITH RECURSIVE children_folders(id, name, folder_id, depth_level) AS(
			SELECT id, name, folder_id, 1 depth_level
			FROM folders
			WHERE folder_id IS NULL
			UNION ALL
			SELECT f.id, f.name, f.folder_id, depth_level + 1
			FROM folders f
			JOIN children_folders
			ON children_folders.id = f.folder_id
		)
		SELECT id, name, folder_id, depth_level
		FROM children_folders
		`).
		Scan(&entities).Error
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting folder tree")
	}

	return entities, nil
}

// GetFoldersAndMedias получение папок и медиа
func (r folderRepository) GetFoldersAndMedias(ctx context.Context, filter actions.FolderFilter) (*actions.FoldersAndMediasList, error) {
	res := &actions.FoldersAndMediasList{
		Total: 0,
		Items: []actions.FolderAndMedia{},
	}

	table := `
		(WITH RECURSIVE tree(id, folder_id, name, updated_at, size) AS (
			SELECT id, folder_id, name, updated_at,
				(SELECT coalesce(SUM(media.size), 0)::bigint FROM media WHERE media.folder_id = id) AS size 
			FROM folders
			UNION ALL
			SELECT fol.id, fol.folder_id, fol.name, fol.updated_at, t.size
			FROM folders fol
			INNER JOIN tree t On fol.id = t.folder_id
		) 
		SELECT id, folder_id, name, filepath, size, updated_at, 'file' AS resource_type
			FROM media
		UNION
		SELECT id, folder_id, name, '', SUM(size), updated_at, 'folder' AS resource_type
			FROM tree
		GROUP BY id, updated_at, folder_id, name) AS fm
		`

	sort := "fm." + filter.Sort
	if filter.Sort == "" {
		sort = "fm.updated_at"
	}
	order := filter.Order
	if filter.Order == "" {
		if filter.Sort == "" {
			order = "desc"
		} else {
			order = "asc"
		}
	}
	err := r.db.Table(table).
		WithContext(ctx).
		Where("fm.folder_id", filter.Filter.FolderId).
		Order(fmt.Sprintf("%s %s", sort, order)).
		Limit(filter.Limit).Offset(filter.Offset).
		Scan(&res.Items).
		Error
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting folders & medias")
	}

	err = r.db.Table(table).
		WithContext(ctx).
		Where("fm.folder_id", filter.Filter.FolderId).
		Count(&res.Total).
		Error
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting folders & medias count")
	}

	return res, nil
}

func (r folderRepository) filterFolder(db *gorm.DB, filter actions.Filter) *gorm.DB {
	if filter.WithFolderId {
		db = db.Where("folder_id", filter.FolderId)
	}
	if filter.Name != "" {
		db = db.Where("name = ?", filter.Name)
	}

	return db
}

func (r folderRepository) filterScopes(filter actions.Filter) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return r.filterFolder(db, filter)
	}
}

// GetFolderParents получение родительских папок
func (r folderRepository) GetFolderParents(ctx context.Context, filter actions.Filter) ([]actions.FolderResponse, error) {
	res := make([]actions.FolderResponse, 0)
	table := "tree"
	err := r.db.WithContext(ctx).
		Table(table).
		Scopes(r.withRecursive(table, filter.FolderId)).
		Order(clause.OrderByColumn{
			Column: clause.Column{Table: table, Name: "depth"},
			Desc:   true,
		}).Scan(&res).Error

	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting parent folders")
	}

	return res, nil
}

func (r folderRepository) withRecursive(tableName string, folderId *uuid.UUID) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Clauses(exclause.With{
			Recursive: true,
			CTEs: []exclause.CTE{{
				Name: tableName,
				Subquery: clause.Expr{
					SQL: `
					SELECT id, folder_id, name, 1 AS depth
					FROM folders
					WHERE id = ?
					UNION ALL
					SELECT f.id, f.folder_id, f.name, depth + 1
					FROM folders f INNER JOIN ? pf
                              ON f.id = pf.folder_id
				 	`,
					Vars: []any{folderId, clause.Table{Name: tableName}},
				},
			}},
		})
	}
}
