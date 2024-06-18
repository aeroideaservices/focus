package handlers

import (
	"github.com/aeroideaservices/focus/page/plugin/actions"
	"github.com/aeroideaservices/focus/page/rest/services"
	middleware "github.com/aeroideaservices/focus/services/gin-middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.aeroidea.ru/internal-projects/focus/services/errors"
	"net/http"
)

type VideoHandler struct {
	videoUseCase *actions.VideoUseCase
	errorHandler *middleware.ErrorHandler
	validator    services.Validator
}

func NewVideoHandler(
	videoUseCase *actions.VideoUseCase,
	errorHandler *middleware.ErrorHandler,
	validator services.Validator,
) *VideoHandler {
	return &VideoHandler{
		videoUseCase: videoUseCase,
		errorHandler: errorHandler,
		validator:    validator,
	}
}

func (h VideoHandler) Create(c *gin.Context) {
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

	resp, err := h.videoUseCase.Create(
		actions.CreateVideoRequest{
			Filename: file.Filename,
			Size:     file.Size,
			FolderId: folderId,
			File:     fo,
		},
	)
	if err != nil {
		_ = c.Error(errors.Internal.Wrap(err, "error creating video"))
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h VideoHandler) GenerateSubtitles(c *gin.Context) {
	var mediaStrIds []string
	if err := c.ShouldBindJSON(&mediaStrIds); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	mediaIds, err := services.GetIdsFromStrings(mediaStrIds)
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	err = h.videoUseCase.GenerateSubtitles(c, mediaIds)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h VideoHandler) UpdateSubtitles(c *gin.Context) {
	mediaId, err := uuid.Parse(c.Params.ByName("media-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	var subtitles actions.SubtitlesToSave

	if err := c.ShouldBindJSON(&subtitles); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.videoUseCase.UpdateSubtitles(c, subtitles, mediaId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}
