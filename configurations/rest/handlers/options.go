package handlers

import (
	"github.com/aeroideaservices/focus/configurations/plugin/actions"
	"github.com/aeroideaservices/focus/configurations/rest/services"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// OptionsHandler обработчик запросов к настройкам
type OptionsHandler struct {
	options   *actions.Options
	validator services.Validator
}

// NewOptionsHandler конструктор
func NewOptionsHandler(
	options *actions.Options,
	validator services.Validator,
) *OptionsHandler {
	return &OptionsHandler{
		options:   options,
		validator: validator,
	}
}

// List получение списка настроек
func (h OptionsHandler) List(c *gin.Context) {
	dto := actions.ListOptions{}
	stringConfId := c.Param(ConfigurationIdParam)
	confId, err := uuid.Parse(stringConfId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}

	err = services.GetLimitAndOffset(c, &dto.Limit, &dto.Offset)
	if err != nil {
		_ = c.Error(err)
		return
	}

	dto.Sort = c.Query("sort")
	dto.Order = c.Query("order")

	dto.ConfId = confId

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	opt, err := h.options.List(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, opt)
}

// Create создание настройки
func (h OptionsHandler) Create(c *gin.Context) {
	stringConfId := c.Param(ConfigurationIdParam)
	confId, err := uuid.Parse(stringConfId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}

	dto := actions.CreateOption{}
	err = c.ShouldBindJSON(&dto)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error parsing json"))
		return
	}

	dto.ConfId = confId

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	optId, err := h.options.Create(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{"id": optId})
}

// UpdateList обновление списка настроек
func (h OptionsHandler) UpdateList(c *gin.Context) {
	var dto actions.UpdateOptionsList
	err := c.ShouldBindJSON(&dto.Items)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error parsing json"))
		return
	}

	stringConfId := c.Param(ConfigurationIdParam)
	confId, err := uuid.Parse(stringConfId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.ConfId = confId

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.options.UpdateList(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Get получение настройки
func (h OptionsHandler) Get(c *gin.Context) {
	stringId := c.Param(OptionIdParam)
	optId, err := uuid.Parse(stringId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}

	stringConfId := c.Param(ConfigurationIdParam)
	confId, err := uuid.Parse(stringConfId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}

	dto := actions.GetOption{Id: optId, ConfId: confId}

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	opt, err := h.options.Get(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, opt)
}

// Update обновление настройки
func (h OptionsHandler) Update(c *gin.Context) {
	stringOptId := c.Param(OptionIdParam)
	optId, err := uuid.Parse(stringOptId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}

	stringConfId := c.Param(ConfigurationIdParam)
	confId, err := uuid.Parse(stringConfId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}

	dto := actions.UpdateOption{}
	err = c.ShouldBindJSON(&dto)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error parsing json"))
		return
	}

	dto.Id = optId
	dto.ConfId = confId

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.options.Update(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Delete удаление настройки
func (h OptionsHandler) Delete(c *gin.Context) {
	stringOptId := c.Param(OptionIdParam)
	optId, err := uuid.Parse(stringOptId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}

	stringConfId := c.Param(ConfigurationIdParam)
	confId, err := uuid.Parse(stringConfId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}

	dto := actions.GetOption{Id: optId, ConfId: confId}

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.options.Delete(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
