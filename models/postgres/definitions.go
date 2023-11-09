package postgres

import (
	"github.com/aeroideaservices/focus/models/plugin/focus"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

var Definitions = []di.Def{
	{
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("focus.db").(*gorm.DB)
			modelsRegistry := ctn.Get("focus.models.registry").(*focus.ModelsRegistry)

			repoResolver := NewRepositoryResolver(db, modelsRegistry)
			return repoResolver, nil
		},
		Name: "focus.models.repositories.resolver",
	},
	{
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("focus.db").(*gorm.DB)
			repo := NewExportInfoRepository(db)
			return repo, nil
		},
		Name: "focus.models.repositories.export",
	},
}
