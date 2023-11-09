package postgres

import (
	"github.com/aeroideaservices/focus/menu/postgres/repositories"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

// Definitions
var Definitions = []di.Def{
	{
		Name: "focus.menu.repositories.menus",
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("focus.db").(*gorm.DB)

			return repositories.NewMenuRepository(db), nil
		},
	},
	{
		Name: "focus.menu.repositories.menuItems",
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("focus.db").(*gorm.DB)

			return repositories.NewMenuItemRepository(db), nil
		},
	},
	{
		Name: "focus.menu.repositories.domains",
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("focus.db").(*gorm.DB)

			return repositories.NewDomainRepository(db), nil
		},
	},
}
