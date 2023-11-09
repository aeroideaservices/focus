package plugin

import (
	"github.com/aeroideaservices/focus/configurations/plugin/actions"
	focsCallbacks "github.com/aeroideaservices/focus/services/callbacks"
	"github.com/sarulabs/di/v2"
)

var Definitions = []di.Def{
	{
		Name: "focus.configurations.actions.configurations",
		Build: func(ctn di.Container) (interface{}, error) {
			confRepository := ctn.Get("focus.configurations.repositories.configurations").(actions.ConfigurationsRepository)

			var callbacks focsCallbacks.Callbacks
			if callbacksI, _ := ctn.SafeGet("focus.configurations.actions.configurations.callbacks"); callbacksI != nil {
				callbacks = callbacksI.(focsCallbacks.Callbacks)
			}

			return actions.NewConfigurations(confRepository, callbacks), nil
		},
	},
	{
		Name: "focus.configurations.actions.options",
		Build: func(ctn di.Container) (interface{}, error) {
			confRepository := ctn.Get("focus.configurations.repositories.configurations").(actions.ConfigurationsRepository)
			optRepository := ctn.Get("focus.configurations.repositories.options").(actions.OptionsRepository)

			var mediaService actions.MediaService
			if mediaServiceInterface, err := ctn.SafeGet("focus.media.actions.media"); err == nil {
				mediaService = mediaServiceInterface.(actions.MediaService)
			}

			var callbacks focsCallbacks.Callbacks
			if callbacksI, _ := ctn.SafeGet("focus.configurations.actions.options.callbacks"); callbacksI != nil {
				callbacks = callbacksI.(focsCallbacks.Callbacks)
			}

			return actions.NewOptions(confRepository, optRepository, mediaService, callbacks), nil
		},
	},
}
