package postgres

import (
	"github.com/aeroideaservices/focus/media/postgres/repositories"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

var Definitions = []di.Def{
	{
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("focus.db").(*gorm.DB)
			return repositories.NewFolderRepository(db), nil
		},
		Name: "focus.media.repository.folder",
	},
	{
		Build: func(ctn di.Container) (interface{}, error) {
			db := ctn.Get("focus.db").(*gorm.DB)
			return repositories.NewMediaRepository(db), nil
		},
		Name: "focus.media.repository.media",
	},
}
