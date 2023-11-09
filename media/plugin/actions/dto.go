package actions

import (
	"github.com/aeroideaservices/focus/media/plugin/service/utils"
	"github.com/google/uuid"
	"io"
)

type CreateFolder struct {
	Name           string     `json:"name" validate:"required,notBlank,min=3,max=50"`
	ParentFolderId *uuid.UUID `json:"parentFolderId,omitempty" validate:"omitempty,notBlank"`
}

type RenameFolder struct {
	Id   uuid.UUID `json:"id" validate:"required,notBlank"`
	Name string    `json:"name" validate:"required,notBlank,min=3,max=50"`
}

type MoveFolder struct {
	Id             uuid.UUID  `json:"id" validate:"required,notBlank"`
	ParentFolderId *uuid.UUID `json:"parentFolderId,omitempty" validate:"omitempty,notBlank"`
}

type FolderAndMediasList struct {
	Total       int64              `json:"total"`
	Items       []FolderAndMedias  `json:"items"`
	Breadcrumbs []FolderBreadcrumb `json:"breadcrumbs"`
}

type FolderBreadcrumb struct {
	Name     string     `json:"name"`
	FolderId *uuid.UUID `json:"folderId"`
}

type FolderAndMedias struct {
	ResourceType string        `json:"resourceType"`
	FolderFields *FolderFields `json:"folderFields,omitempty"`
	FileFields   *FileFields   `json:"fileFields,omitempty"`
}

type FolderFields struct {
	Id   uuid.UUID      `json:"id"`
	Name string         `json:"name"`
	Size utils.Filesize `json:"size"`
}

type FileFields struct {
	Id   uuid.UUID      `json:"id"`
	Name string         `json:"name"`
	Size utils.Filesize `json:"size"`
	Url  string         `json:"url"`
	Ext  string         `json:"ext"`
}

type MediaPreview struct {
	Id          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Ext         string         `json:"ext"`
	Size        utils.Filesize `json:"size"`
	Alt         string         `json:"alt"`
	Title       string         `json:"title"`
	ContentType string         `json:"contentType"`
	Url         string         `json:"url"`
	UpdatedAt   utils.Time     `json:"updatedAt"`
	FolderId    *uuid.UUID     `json:"folderId"`
}

type CreateMedia struct {
	Filename string        `validate:"required,min=3"`
	Size     int64         `validate:""`
	Alt      string        `validate:"omitempty,min=3,max=50"`
	Title    string        `validate:"omitempty,min=3,max=50"`
	FolderId *uuid.UUID    `validate:"omitempty,notBlank"`
	File     io.ReadSeeker `validate:"required"`
}

type CreateMediasList struct {
	FolderId *uuid.UUID  `validate:"omitempty,notBlank"`
	Files    []MediaFile `validate:"required,min=1,max=10,unique=Filename,dive"`
}

type MediaFile struct {
	Filename string        `validate:"required,min=3"`
	Size     int64         `validate:""`
	File     io.ReadSeeker `validate:"required"`
}

type RenameMedia struct {
	Id   uuid.UUID `validate:"required,notBlank"`
	Name string    `validate:"required,notBlank,min=3,max=50" json:"name"`
}

type MoveMedia struct {
	Id       uuid.UUID  `validate:"required,notBlank"`
	FolderId *uuid.UUID `validate:"omitempty,notBlank" json:"folderId"`
}

type MediaShortList struct {
	Items []MediaShort `json:"items"`
}

type MediaShort struct {
	Id    uuid.UUID `json:"id"`
	Alt   string    `json:"alt"`
	Title string    `json:"title"`
	Url   string    `json:"url"`
}

type ListMediasShorts struct {
	Ids []uuid.UUID `json:"ids" validate:"required,min=1,max=1000"`
}

type GetMedia struct {
	Id uuid.UUID `json:"id" validate:"required,notBlank"`
}

type GetFolder struct {
	Id uuid.UUID `json:"id" validate:"required,notBlank"`
}
