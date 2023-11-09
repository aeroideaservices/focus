package handlers

import (
	"github.com/aeroideaservices/focus/models/plugin/actions"
	"github.com/aeroideaservices/focus/models/rest/services"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ModelsHandler обработчик запросов к моделям
type ModelsHandler struct {
	models    *actions.Models
	validator services.Validator
}

// NewModelsHandler конструктор
func NewModelsHandler(
	models *actions.Models,
	validator services.Validator,
) *ModelsHandler {
	return &ModelsHandler{
		models:    models,
		validator: validator,
	}
}

// List получение списка моделей
func (h ModelsHandler) List(c *gin.Context) {
	action := actions.ListModels{
		ModelOrderBy: actions.ModelOrderBy{
			Sort:  c.Query("sort"),
			Order: c.Query("order"),
		},
	}

	var err error
	action.Offset, action.Limit, err = services.GetOffsetAndLimit(c)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	models := h.models.List(action)
	c.JSON(http.StatusOK, models)
}

// Get получение детальной информации о модели
func (h ModelsHandler) Get(c *gin.Context) {
	modelCode := c.Param(ModelCodeParam)
	action := actions.GetModel{
		ModelCode: modelCode,
	}

	err := h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	models, err := h.models.Get(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, models)
}

// GetFieldValues получение значений поля модели
func (h ModelsHandler) GetFieldValues(c *gin.Context) {
	action := actions.ListFieldValues{}
	action.ModelCode = c.Param(ModelCodeParam)
	action.FieldCode = c.Param(FieldCodeParam)
	var err error
	action.Offset, action.Limit, err = services.GetOffsetAndLimit(c)
	if err != nil {
		_ = c.Error(err)
		return
	}
	action.Query = c.Query("query")

	err = h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	fieldValuesList, err := h.models.ListFieldValues(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, fieldValuesList)
}
