package actions

import (
	"context"
	"github.com/aeroideaservices/focus/services/db/db_types/json"
	"io"

	"github.com/google/uuid"

	entity2 "github.com/aeroideaservices/focus/media/plugin/entity"
	"github.com/aeroideaservices/focus/media/plugin/service/utils"
)

type Filter struct {
	Name         string
	FolderId     *uuid.UUID `validate:"omitempty,notBlank"`
	WithFolderId bool
}

type FolderDetail struct {
	Id       uuid.UUID      `json:"id"`
	Name     string         `json:"name"`
	Size     utils.Filesize `json:"size"`
	FolderId *uuid.UUID     `json:"parentFolderId"`
}

type FolderResponse struct {
	Id         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	DepthLevel uint       `json:"depthLevel"`
	FolderId   *uuid.UUID `json:"parentFolderId"`
}
type FolderFilter struct {
	Limit  int    `validate:"required,min=10,max=100"`
	Offset int    `validate:"min=0"`
	Sort   string ``
	Order  string `validate:"omitempty,oneof=asc desc"`
	Filter Filter ``
}
type FoldersAndMediasList struct {
	Total int64
	Items []FolderAndMedia
}

type FolderAndMedia struct {
	ResourceType string         `json:"resourceType"`
	Id           uuid.UUID      `json:"id"`
	Name         string         `json:"name"`
	Size         utils.Filesize `json:"size"`
	FolderId     *uuid.UUID     `json:"parentFolderId"`
	UpdatedAt    utils.Time     `json:"updatedAt"`
	Filepath     string         `json:"filepath,omitempty"`
	Ext          string         `json:"ext,omitempty"`
}

type FolderRepository interface {
	Has(ctx context.Context, id uuid.UUID) bool
	HasByFilter(ctx context.Context, filter Filter) bool

	Get(ctx context.Context, u uuid.UUID) (*entity2.Folder, error)
	GetWithSize(ctx context.Context, u uuid.UUID) (*FolderDetail, error)
	Create(ctx context.Context, folders ...*entity2.Folder) error
	Update(ctx context.Context, folder *entity2.Folder) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter Filter) ([]*entity2.Folder, error)

	HasSubFolder(ctx context.Context, id uuid.UUID, subFolderId *uuid.UUID) (bool, error)

	GetFolderPath(ctx context.Context, id uuid.UUID) (folderPath string, err error)
	GetFolderMediaFilePaths(ctx context.Context, id *uuid.UUID) (mediaFilepath []string, err error)
	GetAllFolderMedias(ctx context.Context, id uuid.UUID) ([]*UpdateMediaDto, error)
	GetFoldersTree(ctx context.Context) ([]*FolderResponse, error)
	GetFoldersAndMedias(ctx context.Context, filter FolderFilter) (*FoldersAndMediasList, error)
	GetFolderParents(ctx context.Context, filter Filter) ([]FolderResponse, error)
}

type UpdateMediaDto struct {
	Id          uuid.UUID  `gorm:"column:id"`
	Name        string     `gorm:"column:name"`
	Filename    string     `gorm:"column:filename"`
	Filepath    string     `gorm:"column:filepath"`
	FolderId    *uuid.UUID `gorm:"column:folder_id"`
	NewFilepath string     `gorm:"column:new_filepath"`
}

type MediaFilter struct {
	FolderId     *uuid.UUID
	WithFolderId bool
	Filename     string
	Filenames    []string
	Name         string
	Ext          string
	Filepath     interface{}
	InIds        []uuid.UUID
}

type MediaRepository interface {
	Has(ctx context.Context, id uuid.UUID) bool
	HasByFilter(ctx context.Context, filter MediaFilter) bool
	HasByFilterWithId(ctx context.Context, filter MediaFilter) (bool, uuid.UUID)
	Create(ctx context.Context, medias ...entity2.Media) error
	Get(ctx context.Context, id uuid.UUID) (*entity2.Media, error)
	Update(ctx context.Context, medias ...*UpdateMediaDto) error
	Delete(ctx context.Context, ids ...uuid.UUID) error
	GetShortList(ctx context.Context, ids []uuid.UUID) ([]entity2.Media, error)
	Count(ctx context.Context, filter MediaFilter) (int, error)

	UpdateSubtitles(ctx context.Context, id uuid.UUID, subtitles json.JSONB) error
}

type UploadFile struct {
	Key         string
	ContentType string
	File        io.ReadSeeker
}

type FileStorage interface {
	Upload(ctx context.Context, media *UploadFile) error
	UploadList(ctx context.Context, media ...UploadFile) error
	Delete(ctx context.Context, keys ...string) error
	Move(ctx context.Context, oldKey string, newKey string) error
	GetSize(ctx context.Context, key string) (int64, error)
	DownloadFile(ctx context.Context, key string, fileName string) error
}

type MediaProvider interface {
	GetUrlByFilepath(mediaFilepath string) string
	GetUrlById(mediaId uuid.UUID) (string, error)
}
