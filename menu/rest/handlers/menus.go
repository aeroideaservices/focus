package handlers

import (
	"github.com/aeroideaservices/focus/menu/plugin/actions"
	"github.com/aeroideaservices/focus/menu/rest/services"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type MenuHandler struct {
	menus     *actions.Menus
	validator services.Validator
}

func NewMenuHandler(
	menus *actions.Menus,
	validator services.Validator,
) *MenuHandler {
	return &MenuHandler{
		menus:     menus,
		validator: validator,
	}
}

func (h MenuHandler) List(c *gin.Context) {
	dto := actions.ListMenus{}

	err := services.GetLimitAndOffset(c, &dto.Limit, &dto.Offset)
	if err != nil {
		_ = c.Error(err)
		return
	}
	dto.Sort = c.Query("sort")
	dto.Order = c.Query("order")

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	menusList, err := h.menus.List(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, menusList)
}

func (h MenuHandler) Create(c *gin.Context) {
	dto := actions.CreateMenu{}
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrapf(err, "json binding error"))
		return
	}

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	menuId, err := h.menus.Create(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{"id": menuId})
}

func (h MenuHandler) Get(c *gin.Context) {
	dto := actions.GetMenu{}

	stringMenuId := c.Param(MenuIdParam)
	menuId, err := uuid.Parse(stringMenuId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.MenuId = menuId

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	menu, err := h.menus.Get(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, menu)
}

func (h MenuHandler) Update(c *gin.Context) {
	dto := actions.UpdateMenu{}
	err := c.ShouldBindJSON(&dto)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "json binding error"))
		return
	}

	stringMenuId := c.Param(MenuIdParam)
	menuId, err := uuid.Parse(stringMenuId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.Id = menuId

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.menus.Update(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h MenuHandler) Delete(c *gin.Context) {
	dto := actions.GetMenu{}

	stringMenuId := c.Param(MenuIdParam)
	menuId, err := uuid.Parse(stringMenuId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.MenuId = menuId

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.menus.Delete(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
