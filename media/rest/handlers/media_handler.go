package handlers

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/aeroideaservices/focus/media/plugin/actions"
	"github.com/aeroideaservices/focus/media/rest/services"
	"github.com/aeroideaservices/focus/services/errors"
)

// MediaHandler обработчик запросов к медиа
type MediaHandler struct {
	medias    *actions.Medias
	validator services.Validator
}

// NewMediaHandler конструктор
func NewMediaHandler(
	medias *actions.Medias,
	validator services.Validator,
) *MediaHandler {
	return &MediaHandler{
		medias:    medias,
		validator: validator,
	}
}

// Upload загрузка медиа
func (h MediaHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error getting form file"))
		return
	}

	fo, err := file.Open()
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error opening form file"))
		return
	}
	defer func() { _ = fo.Close() }()

	var folderId *uuid.UUID
	stringFolderId, hasFolderId := c.GetQuery("folderId")
	if hasFolderId {
		id, err := uuid.Parse(stringFolderId)
		if err != nil {
			_ = c.Error(errors.BadRequest.Wrap(err, "error parsing uuid"))
			return
		}
		folderId = &id
	}

	action := actions.CreateMedia{
		Filename: file.Filename,
		Size:     file.Size,
		FolderId: folderId,
		File:     fo,
	}

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	url, err := h.medias.Upload(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": 1, "file": gin.H{"url": url}})
}

// Create создание медиа
func (h MediaHandler) Create(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error getting form file"))
		return
	}

	fo, err := file.Open()
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error opening form file"))
		return
	}
	defer func() { _ = fo.Close() }()

	var folderId *uuid.UUID
	stringFolderId, hasFolderId := c.GetPostForm("folderId")
	if hasFolderId {
		id, err := uuid.Parse(stringFolderId)
		if err != nil {
			_ = c.Error(errors.BadRequest.Wrap(err, "error parsing uuid"))
			return
		}
		folderId = &id
	}

	action := actions.CreateMedia{
		Filename: file.Filename,
		Size:     file.Size,
		Alt:      c.PostForm("alt"),
		Title:    c.PostForm("title"),
		FolderId: folderId,
		File:     fo,
	}

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	id, err := h.medias.Create(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{"id": id})
}

// UploadList загрузка нескольких медиа
func (h MediaHandler) UploadList(c *gin.Context) {
	action := actions.CreateMediasList{}

	var folderId *uuid.UUID
	stringFolderId, hasFolderId := c.GetPostForm("folderId")
	if hasFolderId {
		id, err := uuid.Parse(stringFolderId)
		if err != nil {
			_ = c.Error(errors.BadRequest.Wrap(err, "error parsing uuid"))
			return
		}
		folderId = &id
	}
	action.FolderId = folderId

	form, err := c.MultipartForm()
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error parsing multipart form"))
		return
	}

	var files []*multipart.FileHeader
	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			files = append(files, fileHeader)
		}
	}
	if len(files) == 0 {
		_ = c.Error(errors.BadRequest.New("files are required"))
		return
	}

	mediaFiles := make([]actions.MediaFile, len(files))
	for i, file := range files {
		fo, err := file.Open()
		if err != nil {
			_ = c.Error(errors.BadRequest.Newf("cannot open file %s", file.Filename))
			return
		}
		defer func(fo multipart.File) { _ = fo.Close() }(fo)
		mediaFiles[i] = actions.MediaFile{
			Filename: file.Filename,
			Size:     file.Size,
			File:     fo,
		}
	}
	action.Files = mediaFiles

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	ids, err := h.medias.UploadList(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{"ids": ids})
}

// Get получение медиа
func (h MediaHandler) Get(c *gin.Context) {
	stringId := c.Param(FileIdParam)
	id, err := uuid.Parse(stringId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error parsing uuid"))
		return
	}

	action := actions.GetMedia{Id: id}

	if err := h.validator.Validate(c, action); err != nil {
		_ = c.Error(err)
		return
	}

	media, err := h.medias.Get(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, media)
}

// Delete удаление медиа
func (h MediaHandler) Delete(c *gin.Context) {
	stringMediaId := c.Param(FileIdParam)
	mediaId, err := uuid.Parse(stringMediaId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error parsing uuid"))
		return
	}

	action := actions.GetMedia{Id: mediaId}

	if err := h.validator.Validate(c, action); err != nil {
		_ = c.Error(err)
		return
	}

	err = h.medias.Delete(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Rename переименование медиа
func (h MediaHandler) Rename(c *gin.Context) {
	stringId := c.Param(FileIdParam)
	mediaId, err := uuid.Parse(stringId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error parsing uuid"))
		return
	}

	action := actions.RenameMedia{}
	err = c.ShouldBindJSON(&action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error parsing json"))
		return
	}
	action.Id = mediaId

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.medias.Rename(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Move перемещение медиа
func (h MediaHandler) Move(c *gin.Context) {
	stringMediaId := c.Param(FileIdParam)
	mediaId, err := uuid.Parse(stringMediaId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error parsing uuid"))
		return
	}

	action := actions.MoveMedia{}
	err = c.ShouldBindJSON(&action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error parsing json"))
		return
	}
	action.Id = mediaId

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.medias.Move(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
