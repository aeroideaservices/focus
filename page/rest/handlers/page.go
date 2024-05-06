package handlers

import (
	"github.com/aeroideaservices/focus/services/errors"
	middleware "github.com/aeroideaservices/focus/services/gin-middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
	"pages/pkg/page/plugin/actions"
	"pages/pkg/page/rest/services"
)

type PageHandler struct {
	pageUseCase  *actions.PageUseCase
	errorHandler *middleware.ErrorHandler
	validator    services.Validator
}

const UniqueViolationErr = "23505"

func NewPageHandler(
	pageUseCase *actions.PageUseCase, errorHandler *middleware.ErrorHandler, validator services.Validator,
) *PageHandler {
	return &PageHandler{
		pageUseCase:  pageUseCase,
		errorHandler: errorHandler,
		validator:    validator,
	}

}

func (h PageHandler) GetById(c *gin.Context) {
	dto := actions.GetPageDto{
		ID: c.Param("page-id"),
	}

	err := h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	page, err := h.pageUseCase.GetById(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if page.ID == uuid.Nil {
		_ = c.Error(errors.NotFound.Wrap(err, "error while getting page"))
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h PageHandler) Create(c *gin.Context) {
	request := &actions.CreatePageRequest{}

	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	err := h.validator.Validate(c, request)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	resp, err := h.pageUseCase.Create(c, request)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == UniqueViolationErr {
				c.JSON(
					http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"error":   "Duplicate key value violates unique constraint",
						"message": "Страница с таким кодом уже существует",
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

func (h PageHandler) Delete(c *gin.Context) {
	pageId, err := uuid.Parse(c.Params.ByName("page-id"))
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	err = h.pageUseCase.Delete(c, pageId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h PageHandler) GetList(c *gin.Context) {
	pages, err := h.pageUseCase.GetList(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, pages)
}

func (h PageHandler) PatchProperties(c *gin.Context) {
	request := &actions.PatchPageRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	pageId, err := uuid.Parse(c.Params.ByName("page-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	request.ID = pageId

	err = h.pageUseCase.PatchProperties(c, request)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == UniqueViolationErr {
				c.JSON(
					http.StatusBadRequest, gin.H{
						"code":    http.StatusBadRequest,
						"error":   "Duplicate key value violates unique constraint",
						"message": "Страница с таким кодом уже существует",
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

func (h PageHandler) PatchGalleryPosition(c *gin.Context) {
	request := &actions.PatchGalleryPosition{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	pageId, err := uuid.Parse(c.Params.ByName("page-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	galleryId, err := uuid.Parse(c.Params.ByName("gallery-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "gallery uuids parsing error"))
		return
	}

	request.PageID = pageId
	request.GalleryID = galleryId

	err = h.validator.Validate(c, request)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.pageUseCase.PatchGalleryPosition(c, request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h PageHandler) LinkGalleries(c *gin.Context) {
	pageId, err := uuid.Parse(c.Params.ByName("page-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	var galleryStrIds []string
	if err = c.ShouldBindJSON(&galleryStrIds); err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "gallery uuids parsing error"))
		return
	}

	galleryIds, err := services.GetIdsFromStrings(galleryStrIds)
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	err, isUniqViolErr := h.pageUseCase.LinkGalleries(c, pageId, galleryIds)
	if err != nil {
		if isUniqViolErr {
			c.JSON(
				http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"error":   "Bad Request",
					"message": "Такая галерея уже есть на странице",
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

func (h PageHandler) UnlinkGalleries(c *gin.Context) {
	pageId, err := uuid.Parse(c.Params.ByName("page-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	var galleryStrIds []string
	if err := c.ShouldBindJSON(&galleryStrIds); err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "gallery uuids parsing error"))
		return
	}

	galleryIds, err := services.GetIdsFromStrings(galleryStrIds)
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "uuid parsing error"))
		return
	}

	err = h.pageUseCase.UnlinkGalleries(c, pageId, galleryIds)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")

}
