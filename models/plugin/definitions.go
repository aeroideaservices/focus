package plugin

import (
	"github.com/aeroideaservices/focus/models/plugin/actions"
	"github.com/aeroideaservices/focus/models/plugin/focus"
	"github.com/aeroideaservices/focus/models/plugin/form"
	focsCallbacks "github.com/aeroideaservices/focus/services/callbacks"
	"github.com/sarulabs/di/v2"
	"go.uber.org/zap"
)

var Definitions = []di.Def{
	{
		Name: "focus.models.actions.models",
		Build: func(ctn di.Container) (interface{}, error) {
			repositoryResolver := ctn.Get("focus.models.repositories.resolver").(actions.RepositoryResolver)
			modelsRegistry := ctn.Get("focus.models.registry").(*focus.ModelsRegistry)
			selectRequest := ctn.Get("focus.models.requests.fieldValues").(form.Request)

			modelsAction := actions.NewModels(modelsRegistry, repositoryResolver, selectRequest)

			return modelsAction, nil
		},
	},
	{
		Name: "focus.models.actions.modelElements",
		Build: func(ctn di.Container) (interface{}, error) {
			repositoryResolver := ctn.Get("focus.models.repositories.resolver").(actions.RepositoryResolver)
			modelsRegistry := ctn.Get("focus.models.registry").(*focus.ModelsRegistry)
			validator := ctn.Get("focus.validator").(actions.Validator)

			var mediaService actions.MediaService
			if mediaServiceI, err := ctn.SafeGet("focus.media.actions.media"); err == nil && mediaServiceI != nil {
				mediaService = mediaServiceI.(actions.MediaService)
			}

			var callbacks = make(map[string]focsCallbacks.Callbacks)
			if callbacksI, _ := ctn.SafeGet("focus.models.actions.modelElements.callbacks"); callbacksI != nil {
				callbacks = callbacksI.(map[string]focsCallbacks.Callbacks)
			}

			modelElementsAction := actions.NewModelElements(modelsRegistry, repositoryResolver, mediaService, validator, callbacks)

			return modelElementsAction, nil
		},
	},
	{
		Name: "focus.models.actions.export",
		Build: func(ctn di.Container) (interface{}, error) {
			repo := ctn.Get("focus.models.repositories.export").(actions.ExportInfoRepository)
			modelsRegistry := ctn.Get("focus.models.registry").(*focus.ModelsRegistry)
			exporter := ctn.Get("focus.models.exporter").(actions.Exporter)
			fileStorage := ctn.Get("focus.models.fileStorage").(actions.FileStorage)
			logger := ctn.Get("focus.logger").(*zap.SugaredLogger)
			fileStorageBaseEndpoint := ctn.Get("focus.models.fileStorage.baseEndpoint").(string)
			return actions.NewExport(repo, modelsRegistry, exporter, fileStorage, logger, fileStorageBaseEndpoint), nil
		},
	},
	{
		Name: "focus.models.registry",
		Build: func(ctn di.Container) (interface{}, error) {
			var supportMedia bool
			models := ctn.Get("focus.models.registry.models").([]any)
			mediaUC, err := ctn.SafeGet("focus.media.actions.media")
			if err == nil && mediaUC != nil {
				supportMedia = true
			}

			modelsRegistry := focus.NewModelsRegistry(supportMedia)
			modelsRegistry.Register(models...)

			return modelsRegistry, nil
		},
	},
}
