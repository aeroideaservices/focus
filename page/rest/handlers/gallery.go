package handlers

import (
	"github.com/aeroideaservices/focus/services/errors"
	middleware "github.com/aeroideaservices/focus/services/gin-middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jemzee04/focus/page/plugin/actions"
	"github.com/jemzee04/focus/page/rest/services"
	"net/http"
)

type GalleryHandler struct {
	galleryUseCase *actions.GalleryUseCase
	errorHandler   *middleware.ErrorHandler
	validator      services.Validator
}

func NewGalleryHandler(
	galleryUseCase *actions.GalleryUseCase, errorHandler *middleware.ErrorHandler, validator services.Validator,
) *GalleryHandler {
	return &GalleryHandler{
		galleryUseCase: galleryUseCase,
		errorHandler:   errorHandler,
		validator:      validator,
	}
}

func (h GalleryHandler) GetList(c *gin.Context) {
	searchValue := c.Query("searchValue")
	pages, err := h.galleryUseCase.GetListWithSearch(c, searchValue)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, pages)
}

func (h GalleryHandler) Create(c *gin.Context) {
	request := &actions.CreateGalleryRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	err := h.validator.Validate(c, request)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	resp, err := h.galleryUseCase.Create(c, request)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == UniqueViolationErr {
				c.JSON(
					http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"error":   "Duplicate key value violates unique constraint",
						"message": "Галерея с таким кодом уже существует",
						"debug":   err.Error(),
					},
				)
				return
			}
		}
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h GalleryHandler) GetById(c *gin.Context) {
	dto := actions.GetGalleryRequest{
		ID: c.Param("gallery-id"),
	}

	err := h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	gallery, err := h.galleryUseCase.GetById(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if gallery.ID == uuid.Nil {
		_ = c.Error(errors.NotFound.Wrap(err, "error while getting gallery"))
		return
	}

	c.JSON(http.StatusOK, gallery)
}

func (h GalleryHandler) Update(c *gin.Context) {
	request := &actions.UpdateGalleryRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	galleryId, err := uuid.Parse(c.Params.ByName("gallery-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "error validating"))
		return
	}

	request.ID = galleryId

	err = h.galleryUseCase.Update(c, request)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == UniqueViolationErr {
				c.JSON(
					http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"error":   "Duplicate key value violates unique constraint",
						"message": "Галерея с таким кодом уже существует",
						"debug":   err.Error(),
					},
				)
				return
			}
		}
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h GalleryHandler) PatchName(c *gin.Context) {
	request := &actions.PatchGalleryNameRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	galleryId, err := uuid.Parse(c.Params.ByName("gallery-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "error validating"))
		return
	}

	request.ID = galleryId

	err = h.validator.Validate(c, request)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.galleryUseCase.PatchName(c, request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h GalleryHandler) PatchCardPosition(c *gin.Context) {
	request := &actions.PatchCardPosition{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	galleryId, err := uuid.Parse(c.Params.ByName("gallery-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "error validating"))
		return
	}

	cardId, err := uuid.Parse(c.Params.ByName("card-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "error validating"))
		return
	}

	request.CardID = cardId
	request.GalleryID = galleryId

	err = h.validator.Validate(c, request)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.galleryUseCase.PatchCardPosition(c, request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h GalleryHandler) LinkCards(c *gin.Context) {
	galleryId, err := uuid.Parse(c.Params.ByName("gallery-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "error validating"))
		return
	}

	var cardStrIds []string
	if err := c.ShouldBindJSON(&cardStrIds); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	cardIds, err := services.GetIdsFromStrings(cardStrIds)
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	err, isUniqViolErr := h.galleryUseCase.LinkCards(c, galleryId, cardIds)
	if err != nil {
		if isUniqViolErr {
			c.JSON(
				http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"error":   "Bad Request",
					"message": "Такая карточка уже есть в галерее",
					"debug":   err.Error(),
				},
			)
			return
		}
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, "success")

}

func (h GalleryHandler) UnlinkCards(c *gin.Context) {
	galleryId, err := uuid.Parse(c.Params.ByName("gallery-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "error validating"))
		return
	}

	var cardStrIds []string
	if err := c.ShouldBindJSON(&cardStrIds); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	cardIds, err := services.GetIdsFromStrings(cardStrIds)
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	err = h.galleryUseCase.UnlinkCards(c, galleryId, cardIds)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")

}

func (h GalleryHandler) DeleteList(c *gin.Context) {
	var galleryStrIds []string
	if err := c.ShouldBindJSON(&galleryStrIds); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	galleryIds, err := services.GetIdsFromStrings(galleryStrIds)
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	err = h.galleryUseCase.DeleteList(c, galleryIds)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}
