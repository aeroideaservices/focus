package postgres

import (
	"fmt"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
	"pages/pkg/page/postgres/repositories"
)

var Definitions = []di.Def{
	{
		Name: "focus.page.repositories.page",
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("db").(*gorm.DB)
			if dialector := db.Dialector.Name(); dialector != "postgres" {
				return nil, fmt.Errorf("focus.page.repositories.page does not support connection %s", dialector)
			}
			return repositories.NewPageRepository(db), nil
		},
	},
	{
		Name: "focus.page.repositories.gallery",
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("focus.db").(*gorm.DB)
			if dialector := db.Dialector.Name(); dialector != "postgres" {
				return nil, fmt.Errorf(
					"focus.gallery.repositories.gallery does not support connection %s", dialector,
				)
			}
			return repositories.NewGalleryRepository(db), nil
		},
	},
	{
		Name: "focus.page.repositories.tag",
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("focus.db").(*gorm.DB)
			if dialector := db.Dialector.Name(); dialector != "postgres" {
				return nil, fmt.Errorf(
					"focus.gallery.repositories.gallery does not support connection %s", dialector,
				)
			}
			return repositories.NewTagRepository(db), nil
		},
	},
	{
		Name: "focus.card.repositories.card",
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("db").(*gorm.DB)
			if dialector := db.Dialector.Name(); dialector != "postgres" {
				return nil, fmt.Errorf(
					"focus.configurations.repositories.configurations does not support connection %s", dialector,
				)
			}
			return repositories.NewCardRepository(db), nil
		},
	},
}
