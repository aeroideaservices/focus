package rest

import (
	middleware "github.com/aeroideaservices/focus/services/gin-middleware"
	"github.com/sarulabs/di/v2"
	"pages/pkg/page/plugin/actions"
	"pages/pkg/page/rest/handlers"
	"pages/pkg/page/rest/services"
)

var Definitions = []di.Def{
	{
		Name: "focus.page.handlers.page",
		Build: func(ctn di.Container) (interface{}, error) {
			pageUseCase := ctn.Get("focus.page.actions.page").(*actions.PageUseCase)
			validator := ctn.Get("focus.validator").(services.Validator)
			errorHandler := ctn.Get("focus.errorHandler").(*middleware.ErrorHandler)
			return handlers.NewPageHandler(pageUseCase, errorHandler, validator), nil
		},
	},
	{
		Name: "focus.page.handlers.gallery",
		Build: func(ctn di.Container) (interface{}, error) {
			galleryUseCase := ctn.Get("focus.page.actions.gallery").(*actions.GalleryUseCase)
			validator := ctn.Get("focus.validator").(services.Validator)
			errorHandler := ctn.Get("focus.errorHandler").(*middleware.ErrorHandler)
			return handlers.NewGalleryHandler(galleryUseCase, errorHandler, validator), nil
		},
	},
	{
		Name: "focus.page.handlers.card",
		Build: func(ctn di.Container) (interface{}, error) {
			cardUseCase := ctn.Get("focus.page.actions.card").(*actions.CardUseCase)
			validator := ctn.Get("focus.validator").(services.Validator)
			errorHandler := ctn.Get("focus.errorHandler").(*middleware.ErrorHandler)
			return handlers.NewCardHandler(cardUseCase, errorHandler, validator), nil
		},
	},
	{
		Name: "focus.page.handlers.tag",
		Build: func(ctn di.Container) (interface{}, error) {
			tagUseCase := ctn.Get("focus.page.actions.tag").(*actions.TagUseCase)
			validator := ctn.Get("focus.validator").(services.Validator)
			errorHandler := ctn.Get("focus.errorHandler").(*middleware.ErrorHandler)
			return handlers.NewTagHandler(tagUseCase, errorHandler, validator), nil
		},
	},
	{
		Name: "focus.page.handlers.video",
		Build: func(ctn di.Container) (interface{}, error) {
			videoUseCase := ctn.Get("focus.page.actions.video").(*actions.VideoUseCase)
			validator := ctn.Get("focus.validator").(services.Validator)
			errorHandler := ctn.Get("focus.errorHandler").(*middleware.ErrorHandler)
			return handlers.NewVideoHandler(videoUseCase, errorHandler, validator), nil
		},
	},
	{
		Name: "focus.page.router",
		Build: func(ctn di.Container) (interface{}, error) {
			pageHandler := ctn.Get("focus.page.handlers.page").(*handlers.PageHandler)
			galleryHandler := ctn.Get("focus.page.handlers.gallery").(*handlers.GalleryHandler)
			cardHandler := ctn.Get("focus.page.handlers.card").(*handlers.CardHandler)
			tagHandleer := ctn.Get("focus.page.handlers.tag").(*handlers.TagHandler)
			//optHandler := ctn.Get("focus.configurations.handlers.options").(*handlers.)
			videoHandler := ctn.Get("focus.page.handlers.video").(*handlers.VideoHandler)
			errorHandler := ctn.Get("focus.errorHandler").(services.ErrorHandler)
			return NewRouter(pageHandler, galleryHandler, cardHandler, tagHandleer, videoHandler, errorHandler), nil
		},
	},
}
