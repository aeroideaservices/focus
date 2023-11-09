package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/aeroideaservices/focus/media/plugin/actions"
	"github.com/aeroideaservices/focus/media/rest/services"
)

// FolderHandler обработчик запросов к папкам
type FolderHandler struct {
	folders   *actions.Folders
	validator services.Validator
}

// NewFolderHandler конструктор
func NewFolderHandler(
	folders *actions.Folders,
	validator services.Validator,
) *FolderHandler {
	return &FolderHandler{
		folders:   folders,
		validator: validator,
	}
}

// GetAll получение папок и медиа в папке
func (h FolderHandler) GetAll(c *gin.Context) {
	action := actions.FolderFilter{}
	var folderId *uuid.UUID
	stringId, hasId := c.GetQuery("parentFolderId")
	if hasId {
		id, err := uuid.Parse(stringId)
		if err != nil {
			_ = c.Error(err)
			return
		}
		folderId = &id
	}

	err := services.GetLimitAndOffset(c, &action.Limit, &action.Offset)
	if err != nil {
		_ = c.Error(err)
		return
	}
	action.Sort = c.Query("sort")
	action.Order = c.Query("order")
	action.Filter = actions.Filter{
		FolderId: folderId,
	}

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	list, err := h.folders.GetAll(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, list)
}

// GetTree получение дерева папок
func (h FolderHandler) GetTree(c *gin.Context) {
	folders, err := h.folders.GetTree(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, folders)
}

// Create создание папки
func (h FolderHandler) Create(c *gin.Context) {
	action := actions.CreateFolder{}
	err := c.ShouldBindJSON(&action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	id, err := h.folders.Create(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{"id": id})
}

// Get получение папки
func (h FolderHandler) Get(c *gin.Context) {
	action := actions.GetFolder{}
	stringId := c.Param(FolderIdParam)
	id, err := uuid.Parse(stringId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	action.Id = id

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	folder, err := h.folders.Get(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, folder)
}

// Delete удаление папки
func (h FolderHandler) Delete(c *gin.Context) {
	action := actions.GetFolder{}
	stringId := c.Param(FolderIdParam)
	id, err := uuid.Parse(stringId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	action.Id = id

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.folders.Delete(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Rename переименование папки
func (h FolderHandler) Rename(c *gin.Context) {
	stringFolderId := c.Param(FolderIdParam)
	folderId, err := uuid.Parse(stringFolderId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	action := actions.RenameFolder{}
	err = c.ShouldBindJSON(&action)
	if err != nil {
		_ = c.Error(err)
		return
	}
	action.Id = folderId

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.folders.Rename(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Move перемещение папки
func (h FolderHandler) Move(c *gin.Context) {
	stringFolderId := c.Param(FolderIdParam)
	folderId, err := uuid.Parse(stringFolderId)
	if err != nil {
		_ = c.Error(err)
		return
	}

	action := actions.MoveFolder{}
	err = c.ShouldBindJSON(&action)
	if err != nil {
		_ = c.Error(err)
		return
	}
	action.Id = folderId

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.folders.Move(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
