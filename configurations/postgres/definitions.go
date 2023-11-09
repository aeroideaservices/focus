package postgres

import (
	"fmt"
	"github.com/aeroideaservices/focus/configurations/postgres/repositories"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

var Definitions = []di.Def{
	{
		Name: "focus.configurations.repositories.configurations",
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("focus.db").(*gorm.DB)
			if dialector := db.Dialector.Name(); dialector != "postgres" {
				return nil, fmt.Errorf("focus.configurations.repositories.configurations does not support connection %s", dialector)
			}
			return repositories.NewConfigurationRepository(db), nil
		},
	},
	{
		Name: "focus.configurations.repositories.options",
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("focus.db").(*gorm.DB)
			if dialector := db.Dialector.Name(); dialector != "postgres" {
				return nil, fmt.Errorf("focus.configurations.repositories.options does not support connection %s", dialector)
			}
			return repositories.NewOptionsRepository(db), nil
		},
	},
}
