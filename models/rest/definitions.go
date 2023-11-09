package rest

import (
	"github.com/aeroideaservices/focus/models/plugin/actions"
	"github.com/aeroideaservices/focus/models/plugin/form"
	"github.com/aeroideaservices/focus/models/rest/handlers"
	"github.com/aeroideaservices/focus/models/rest/services"
	"github.com/sarulabs/di/v2"
	"net/http"
)

// Definitions определение сервисов для контейнера
var Definitions = []di.Def{
	{
		Name: "focus.models.handler.models",
		Build: func(ctn di.Container) (interface{}, error) {
			modelsAction := ctn.Get("focus.models.actions.models").(*actions.Models)
			validator := ctn.Get("focus.validator").(services.Validator)
			return handlers.NewModelsHandler(modelsAction, validator), nil
		},
	},
	{
		Name: "focus.models.handler.modelElements",
		Build: func(ctn di.Container) (interface{}, error) {
			modelElementsAction := ctn.Get("focus.models.actions.modelElements").(*actions.ModelElements)
			validator := ctn.Get("focus.validator").(services.Validator)
			return handlers.NewElementsHandler(modelElementsAction, validator), nil
		},
	},
	{
		Name: "focus.models.handler.export",
		Build: func(ctn di.Container) (interface{}, error) {
			modelExportAction := ctn.Get("focus.models.actions.export").(*actions.Export)
			validator := ctn.Get("focus.validator").(services.Validator)
			return handlers.NewExportHandler(modelExportAction, validator), nil
		},
	},
	{
		Name: "focus.models.router",
		Build: func(ctn di.Container) (interface{}, error) {
			modelsHandler := ctn.Get("focus.models.handler.models").(*handlers.ModelsHandler)
			modelElementsHandler := ctn.Get("focus.models.handler.modelElements").(*handlers.ElementsHandler)
			modelExportHandler := ctn.Get("focus.models.handler.export").(*handlers.ExportHandler)
			errorHandler := ctn.Get("focus.errorHandler").(services.ErrorHandler)
			return NewRouter(modelsHandler, modelElementsHandler, modelExportHandler, errorHandler), nil
		},
	},
	{
		Name: "focus.models.requests.fieldValues",
		Build: func(ctn di.Container) (interface{}, error) {
			return form.Request{
				URI:       "/models-v2/{" + handlers.ModelCodeParam + "}/fields/{" + handlers.FieldCodeParam + "}",
				Meth:      http.MethodGet,
				Paginated: true,
			}, nil
		},
	},
}
