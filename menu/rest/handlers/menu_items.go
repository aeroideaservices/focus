package handlers

import (
	"github.com/aeroideaservices/focus/menu/plugin/actions"
	"github.com/aeroideaservices/focus/menu/rest/services"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type MenuItemHandler struct {
	menuItems *actions.MenuItems
	validator services.Validator
}

func NewMenuItemHandler(
	menuItems *actions.MenuItems,
	validator services.Validator,
) *MenuItemHandler {
	return &MenuItemHandler{
		menuItems: menuItems,
		validator: validator,
	}
}

func (h MenuItemHandler) List(c *gin.Context) {
	dto := actions.ListMenuItems{}

	stringMenuId := c.Param(MenuIdParam)
	menuId, err := uuid.Parse(stringMenuId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.Filter.MenuId = menuId
	dto.Sort = c.Query("sort")
	dto.Order = c.Query("order")
	parentMenuItemIdString := c.Query("parentMenuItemId")
	if parentMenuItemIdString != "" {
		parentMenuItemId, err := uuid.Parse(parentMenuItemIdString)
		if err != nil {
			_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
			return
		}
		dto.Filter.ParentId = &parentMenuItemId
	}

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	menuItemsList, err := h.menuItems.List(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, menuItemsList)
}

func (h MenuItemHandler) Create(c *gin.Context) {
	dto := actions.CreateMenuItem{}

	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "json binding error"))
		return
	}

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

	menuItemId, err := h.menuItems.Create(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{"id": menuItemId})
}

func (h MenuItemHandler) Get(c *gin.Context) {
	dto := actions.GetMenuItem{}

	stringMenuId := c.Param(MenuIdParam)
	menuId, err := uuid.Parse(stringMenuId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.MenuId = menuId

	stringMenuItemId := c.Param(MenuItemIdParam)
	menuItemId, err := uuid.Parse(stringMenuItemId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.MenuItemId = menuItemId

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	menuItem, err := h.menuItems.Get(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, menuItem)
}

func (h MenuItemHandler) Update(c *gin.Context) {
	dto := actions.UpdateMenuItem{}

	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "json binding error"))
		return
	}

	stringMenuId := c.Param(MenuIdParam)
	menuId, err := uuid.Parse(stringMenuId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.MenuId = menuId

	stringMenuItemId := c.Param(MenuItemIdParam)
	menuItemId, err := uuid.Parse(stringMenuItemId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.MenuItemId = menuItemId

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.menuItems.Update(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h MenuItemHandler) Delete(c *gin.Context) {
	dto := actions.GetMenuItem{}

	stringMenuId := c.Param(MenuIdParam)
	menuId, err := uuid.Parse(stringMenuId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.MenuId = menuId

	stringMenuItemId := c.Param(MenuItemIdParam)
	menuItemId, err := uuid.Parse(stringMenuItemId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.MenuItemId = menuItemId

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.menuItems.Delete(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h MenuItemHandler) Move(c *gin.Context) {
	dto := actions.MoveMenuItem{}

	if err := c.ShouldBindJSON(&dto); err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "json binding error"))
		return
	}

	stringMenuId := c.Param(MenuIdParam)
	menuId, err := uuid.Parse(stringMenuId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.MenuId = menuId

	stringMenuItemId := c.Param(MenuItemIdParam)
	menuItemId, err := uuid.Parse(stringMenuItemId)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "uuid parsing error"))
		return
	}
	dto.MenuItemId = menuItemId

	err = h.validator.Validate(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = h.menuItems.Move(c, dto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
