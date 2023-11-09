package actions

import (
	"github.com/aeroideaservices/focus/services/errors"
)

var (
	ErrMenuAlreadyExists = errors.Conflict.New("menu with the same code already exists").T("menu.conflict")
	ErrMaxPosition       = errors.BadRequest.New("position of menu item too large").T("menu-item.position-too-large")
	ErrFieldNotUpdatable = errors.BadRequest.New("field is not updatable")
	ErrMaxDepthExceeded  = errors.BadRequest.New("the maximum depth level for menu items has been reached").T("menu-item.max-depth-exceeded")
	ErrMenuNotFound      = errors.NotFound.New("menu not found").T("menu.not-found")
	ErrMenuItemNotFound  = errors.NotFound.New("menu item does not exist").T("menu-item.not-found")
)
