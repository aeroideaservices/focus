package handlers

import (
	"github.com/aeroideaservices/focus/models/plugin/actions"
	"github.com/aeroideaservices/focus/models/rest/services"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// ElementsHandler обработчик запросов к элементам модели
type ElementsHandler struct {
	elements  *actions.ModelElements
	validator services.Validator
}

// NewElementsHandler конструктор
func NewElementsHandler(
	elements *actions.ModelElements,
	validator services.Validator,
) *ElementsHandler {
	return &ElementsHandler{
		elements:  elements,
		validator: validator,
	}
}

// List получение списка элементов модели
func (h ElementsHandler) List(c *gin.Context) {
	action := actions.ListModelElements{}
	err := c.ShouldBindJSON(&action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	action.Offset, action.Limit, err = services.GetOffsetAndLimit(c)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error getting limit and offset"))
		return
	}

	action.ModelCode = c.Param(ModelCodeParam)
	action.Sort = c.Query("sort")
	action.Order = c.Query("order")
	action.Query = c.Query("query")

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	res, err := h.elements.List(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get получение элемента модели
func (h ElementsHandler) Get(c *gin.Context) {
	action := actions.GetModelElement{
		ModelCode: c.Param(ModelCodeParam),
		PKey:      c.Param(ModelElementIDParam),
	}

	err := h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	res, err := h.elements.Get(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// Create создание нового элемента модели
func (h ElementsHandler) Create(c *gin.Context) {
	action := actions.CreateModelElement{}
	err := c.ShouldBindJSON(&action.ModelElement)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	action.ModelCode = c.Param(ModelCodeParam)

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	id, err := h.elements.Create(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{"id": id})
}

// Update обновление элемента модели
func (h ElementsHandler) Update(c *gin.Context) {
	action := actions.UpdateModelElement{}
	err := c.ShouldBindJSON(&action.ModelElement)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}

	stringId := c.Param(ModelElementIDParam)
	id, err := uuid.Parse(stringId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}
	action.PKey = id
	action.ModelCode = c.Param(ModelCodeParam)

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.elements.Update(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Delete удаление элемента модели
func (h ElementsHandler) Delete(c *gin.Context) {
	action := actions.DeleteModelElement{
		ModelCode: c.Param(ModelCodeParam),
		PKey:      c.Param(ModelElementIDParam),
	}

	err := h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.elements.Delete(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// DeleteList удаление нескольких элементов модели
func (h ElementsHandler) DeleteList(c *gin.Context) {
	action := actions.DeleteModelElements{}

	err := c.ShouldBindJSON(&action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error converting request body to action"))
		return
	}
	action.ModelCode = c.Param(ModelCodeParam)

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.elements.DeleteList(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
