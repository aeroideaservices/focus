package plugin

import (
	"github.com/aeroideaservices/focus/menu/plugin/actions"
	focsCallbacks "github.com/aeroideaservices/focus/services/callbacks"
	"github.com/sarulabs/di/v2"
)

// Definitions
var Definitions = []di.Def{
	{
		Name: "focus.menu.actions.menus",
		Build: func(ctn di.Container) (interface{}, error) {
			menuRepository := ctn.Get("focus.menu.repositories.menus").(actions.MenuRepository)

			var callbacks focsCallbacks.Callbacks
			if callbacksI, _ := ctn.SafeGet("focus.menu.actions.menus.callbacks"); callbacksI != nil {
				callbacks = callbacksI.(focsCallbacks.Callbacks)
			}

			return actions.NewMenus(menuRepository, callbacks), nil
		},
	},
	{
		Name: "focus.menu.actions.menuItems",
		Build: func(ctn di.Container) (interface{}, error) {
			menuRepository := ctn.Get("focus.menu.repositories.menus").(actions.MenuRepository)
			menuItemRepository := ctn.Get("focus.menu.repositories.menuItems").(actions.MenuItemRepository)
			maxMenuItemsDepth := ctn.Get("focus.menu.maxMenuItemsDepth").(int)

			var callbacks focsCallbacks.Callbacks
			if callbacksI, _ := ctn.SafeGet("focus.menu.actions.menuItems.callbacks"); callbacksI != nil {
				callbacks = callbacksI.(focsCallbacks.Callbacks)
			}

			return actions.NewMenuItems(menuRepository, menuItemRepository, maxMenuItemsDepth, callbacks), nil
		},
	},
	{
		Name: "focus.menu.actions.domains",
		Build: func(ctn di.Container) (interface{}, error) {
			domainsRepository := ctn.Get("focus.menu.repositories.domains").(actions.DomainsRepository)

			return actions.NewDomains(domainsRepository), nil
		},
	},
}
