package rest

import (
	"github.com/aeroideaservices/focus/models/rest/handlers"
	"github.com/aeroideaservices/focus/models/rest/services"
	"github.com/gin-gonic/gin"
)

// Router роутер
type Router struct {
	modelsHandler        *handlers.ModelsHandler
	errorHandler         services.ErrorHandler
	modelElementsHandler *handlers.ElementsHandler
	modelExportHandler   *handlers.ExportHandler
}

// NewRouter конструктор
func NewRouter(
	modelsHandler *handlers.ModelsHandler,
	modelElementsHandler *handlers.ElementsHandler,
	modelExportHandler *handlers.ExportHandler,
	errorHandler services.ErrorHandler,
) *Router {
	return &Router{
		modelsHandler:        modelsHandler,
		modelElementsHandler: modelElementsHandler,
		modelExportHandler:   modelExportHandler,
		errorHandler:         errorHandler,
	}
}

// SetRoutes проставление роутов
func (r *Router) SetRoutes(routerGroup *gin.RouterGroup) {
	models := routerGroup.Group("models-v2")
	models.Use(r.errorHandler.Handle) // отлов ошибок

	models.GET("", r.modelsHandler.List)

	model := models.Group(":" + handlers.ModelCodeParam)
	model.GET("", r.modelsHandler.Get)
	//model.GET("/settings", r.modelHandler.GetSettings) todo

	export := model.Group("export")
	export.POST("", r.modelExportHandler.Export)
	export.GET("", r.modelExportHandler.GetExportInfo)

	modelFields := model.Group("fields")
	modelFieldValues := modelFields.Group(":" + handlers.FieldCodeParam)
	modelFieldValues.GET("", r.modelsHandler.GetFieldValues)

	modelElements := model.Group("elements")
	modelElements.POST("", r.modelElementsHandler.Create)
	modelElements.POST("list", r.modelElementsHandler.List)
	modelElements.POST("batch-delete", r.modelElementsHandler.DeleteList)

	modelElement := modelElements.Group(":" + handlers.ModelElementIDParam)
	modelElement.GET("", r.modelElementsHandler.Get)
	modelElement.PUT("", r.modelElementsHandler.Update)
	modelElement.DELETE("", r.modelElementsHandler.Delete)
}
