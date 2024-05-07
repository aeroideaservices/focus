package handlers

import (
	"github.com/aeroideaservices/focus/services/errors"
	middleware "github.com/aeroideaservices/focus/services/gin-middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jemzee04/focus/page/plugin/actions"
	"github.com/jemzee04/focus/page/rest/services"
	"net/http"
)

type TagHandler struct {
	tagUseCase   *actions.TagUseCase
	errorHandler *middleware.ErrorHandler
	validator    services.Validator
}

func NewTagHandler(
	tagUseCase *actions.TagUseCase, errorHandler *middleware.ErrorHandler, validator services.Validator,
) *TagHandler {
	return &TagHandler{
		tagUseCase:   tagUseCase,
		errorHandler: errorHandler,
		validator:    validator,
	}
}

func (h TagHandler) GetList(c *gin.Context) {
	searchValue := c.Query("searchValue")
	pages, err := h.tagUseCase.GetListWithSearch(c, searchValue)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, pages)
}

func (h TagHandler) Create(c *gin.Context) {
	request := &actions.TagDto{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	err := h.validator.Validate(c, request)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	resp, err := h.tagUseCase.Create(c, request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h TagHandler) Update(c *gin.Context) {
	request := &actions.TagDtoRequest{}
	if err := c.ShouldBindJSON(request); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	tagId, err := uuid.Parse(c.Params.ByName("tag-id"))
	if err != nil {
		_ = c.Error(errors.NotFound.Wrap(err, "error validating"))
		return
	}

	request.ID = tagId

	err = h.tagUseCase.Update(c, request)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}

func (h TagHandler) Delete(c *gin.Context) {
	tagId, err := uuid.Parse(c.Params.ByName("tag-id"))
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.tagUseCase.Delete(c, tagId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, "success")
}
