package handlers

import (
	"github.com/aeroideaservices/focus/configurations/plugin/actions"
	"github.com/aeroideaservices/focus/configurations/rest/services"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// ConfigurationsHandler обработчик запросов к конфигурациям
type ConfigurationsHandler struct {
	configurations *actions.Configurations
	validator      services.Validator
}

// NewConfigurationsHandler конструктор
func NewConfigurationsHandler(
	configurations *actions.Configurations,
	validator services.Validator,
) *ConfigurationsHandler {
	return &ConfigurationsHandler{
		configurations: configurations,
		validator:      validator,
	}
}

// List получение списка конфигураций
func (h ConfigurationsHandler) List(c *gin.Context) {
	dto := actions.ListConfigurations{}

	err := services.GetLimitAndOffset(c, &dto.Limit, &dto.Offset)
	if err != nil {
		_ = c.Error(err)
		return
	}

	dto.Sort = c.Query("sort")
	dto.Order = c.Query("order")
	dto.Query = c.Query("query")

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	conf, err := h.configurations.List(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, conf)
}

// Create создание конфигурации
func (h ConfigurationsHandler) Create(c *gin.Context) {
	dto := actions.CreateConfiguration{}
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	confId, err := h.configurations.Create(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{"id": confId})
}

// Get получение конфигурации
func (h ConfigurationsHandler) Get(c *gin.Context) {
	action := actions.GetConfiguration{}
	stringId := c.Param(ConfigurationIdParam)
	id, err := uuid.Parse(stringId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	action.Id = id

	conf, err := h.configurations.Get(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, conf)
}

// Update обновление конфигурации
func (h ConfigurationsHandler) Update(c *gin.Context) {
	dto := actions.UpdateConfiguration{}

	err := c.ShouldBindJSON(&dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	stringId := c.Param(ConfigurationIdParam)
	dto.Id, err = uuid.Parse(stringId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.configurations.Update(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Delete удаление конфигурации
func (h ConfigurationsHandler) Delete(c *gin.Context) {
	action := actions.GetConfiguration{}
	stringId := c.Param(ConfigurationIdParam)
	id, err := uuid.Parse(stringId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	action.Id = id

	err = h.configurations.Delete(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
