package rest

import (
	"github.com/aeroideaservices/focus/menu/plugin/actions"
	"github.com/aeroideaservices/focus/menu/rest/handlers"
	"github.com/aeroideaservices/focus/menu/rest/services"
	"github.com/sarulabs/di/v2"
)

var Definitions = []di.Def{
	{
		Name: "focus.menu.handlers.menus",
		Build: func(ctn di.Container) (interface{}, error) {
			menus := ctn.Get("focus.menu.actions.menus").(*actions.Menus)
			validator := ctn.Get("focus.validator").(services.Validator)
			return handlers.NewMenuHandler(menus, validator), nil
		},
	},
	{
		Name: "focus.menu.handlers.menuItems",
		Build: func(ctn di.Container) (interface{}, error) {
			menuItems := ctn.Get("focus.menu.actions.menuItems").(*actions.MenuItems)
			validator := ctn.Get("focus.validator").(services.Validator)
			return handlers.NewMenuItemHandler(menuItems, validator), nil
		},
	},
	{
		Name: "focus.menu.handlers.domains",
		Build: func(ctn di.Container) (interface{}, error) {
			domains := ctn.Get("focus.menu.actions.domains").(*actions.Domains)
			validator := ctn.Get("focus.validator").(services.Validator)
			return handlers.NewDomainsHandler(domains, validator), nil
		},
	},
	{
		Name: "focus.menu.router",
		Build: func(ctn di.Container) (interface{}, error) {
			menuHandler := ctn.Get("focus.menu.handlers.menus").(*handlers.MenuHandler)
			menuItemHandler := ctn.Get("focus.menu.handlers.menuItems").(*handlers.MenuItemHandler)
			domainsHandler := ctn.Get("focus.menu.handlers.domains").(*handlers.DomainsHandler)
			errorHandler := ctn.Get("focus.errorHandler").(services.ErrorHandler)
			return NewRouter(menuHandler, menuItemHandler, domainsHandler, errorHandler), nil
		},
	},
}
