package rest

import (
	"github.com/aeroideaservices/focus/media/plugin/actions"
	"github.com/sarulabs/di/v2"

	"github.com/aeroideaservices/focus/media/rest/handlers"
	"github.com/aeroideaservices/focus/media/rest/services"
)

var Definitions = []di.Def{
	{
		Build: func(ctn di.Container) (interface{}, error) {
			folders := ctn.Get("focus.media.actions.folder").(*actions.Folders)
			validator := ctn.Get("focus.validator").(services.Validator)
			return handlers.NewFolderHandler(folders, validator), nil
		},
		Name: "focus.media.handler.folder",
	},
	{
		Build: func(ctn di.Container) (interface{}, error) {
			medias := ctn.Get("focus.media.actions.media").(*actions.Medias)
			validator := ctn.Get("focus.validator").(services.Validator)
			return handlers.NewMediaHandler(medias, validator), nil
		},
		Name: "focus.media.handler.media",
	},
	{
		Build: func(ctn di.Container) (interface{}, error) {
			confHandler := ctn.Get("focus.media.handler.folder").(*handlers.FolderHandler)
			optHandler := ctn.Get("focus.media.handler.media").(*handlers.MediaHandler)
			errorHandler := ctn.Get("focus.errorHandler").(services.ErrorHandler)
			return NewRouter(confHandler, optHandler, errorHandler), nil
		},
		Name: "focus.media.router",
	},
}
