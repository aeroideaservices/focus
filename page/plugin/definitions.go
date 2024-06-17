package plugin

import (
	media_usecase "github.com/aeroideaservices/focus/media/plugin/actions"
	"github.com/aeroideaservices/focus/page/plugin/actions"
	"github.com/aeroideaservices/focus/page/plugin/services"
	"github.com/sarulabs/di/v2"
	actions3 "gitlab.aeroidea.ru/internal-projects/focus/forms/plugin/actions"
	//	helpers "gitlab.aeroidea.ru/platform/platformlib/go/lib/golang-helpers-lib"

	"go.uber.org/zap"
)

var Definitions = []di.Def{
	{
		Name: "copier_service",
		Build: func(ctn di.Container) (interface{}, error) {
			return services.Copier{}, nil
		},
	},
	{
		Name: "focus.page.actions.page",
		Build: func(ctn di.Container) (interface{}, error) {
			pageRepository := ctn.Get("focus.page.repositories.page").(actions.PageRepository)
			galleryRepository := ctn.Get("focus.page.repositories.gallery").(actions.GalleryRepository)
			galleryUseCase := ctn.Get("focus.page.actions.gallery").(*actions.GalleryUseCase)
			copierService := ctn.Get("copier_service").(actions.CopierInterface)
			logger := ctn.Get("logger").(*zap.SugaredLogger)
			//var callbacks callbacks.Callbacks
			//if callbacksI, _ := ctn.SafeGet("focus.configurations.actions.configurations.callbacks"); callbacksI != nil {
			//	callbacks = callbacksI.(focsCallbacks.Callbacks)
			//}

			return actions.NewPageUseCase(
				pageRepository, galleryRepository, *galleryUseCase, copierService, logger,
			), nil
		},
	},
	{
		Name: "focus.page.actions.gallery",
		Build: func(ctn di.Container) (interface{}, error) {
			galleryRepository := ctn.Get("focus.page.repositories.gallery").(actions.GalleryRepository)
			cardRepository := ctn.Get("focus.card.repositories.card").(actions.CardRepository)
			cardUseCase := ctn.Get("focus.page.actions.card").(*actions.CardUseCase)
			copierService := ctn.Get("copier_service").(actions.CopierInterface)
			logger := ctn.Get("logger").(*zap.SugaredLogger)
			//var callbacks focsCallbacks.Callbacks
			//if callbacksI, _ := ctn.SafeGet("focus.configurations.actions.configurations.callbacks"); callbacksI != nil {
			//	callbacks = callbacksI.(focsCallbacks.Callbacks)
			//}

			return actions.NewGalleryUseCase(
				galleryRepository, cardRepository, *cardUseCase, copierService, logger,
			), nil
		},
	},
	{
		Name: "focus.page.actions.card",
		Build: func(ctn di.Container) (interface{}, error) {
			cardRepository := ctn.Get("focus.card.repositories.card").(actions.CardRepository)
			galleryRepository := ctn.Get("focus.page.repositories.gallery").(actions.GalleryRepository)
			mediaProvider := ctn.Get("focus.media.provider").(media_usecase.MediaProvider)
			tagRepository := ctn.Get("focus.page.repositories.tag").(actions.TagRepository)
			copierService := ctn.Get("copier_service").(actions.CopierInterface)
			formUseCase := ctn.Get("focus.forms.actions.forms").(*actions3.Forms)
			logger := ctn.Get("logger").(*zap.SugaredLogger)
			//var callbacks focsCallbacks.Callbacks
			//if callbacksI, _ := ctn.SafeGet("focus.configurations.actions.configurations.callbacks"); callbacksI != nil {
			//	callbacks = callbacksI.(focsCallbacks.Callbacks)
			//}

			return actions.NewCardUseCase(
				cardRepository, galleryRepository, tagRepository, mediaProvider, *formUseCase, copierService, logger,
			), nil
		},
	},
	{
		Name: "focus.page.actions.tag",
		Build: func(ctn di.Container) (interface{}, error) {
			tagRepository := ctn.Get("focus.page.repositories.tag").(actions.TagRepository)
			copierService := ctn.Get("copier_service").(actions.CopierInterface)
			logger := ctn.Get("logger").(*zap.SugaredLogger)
			//var callbacks focsCallbacks.Callbacks
			//if callbacksI, _ := ctn.SafeGet("focus.configurations.actions.configurations.callbacks"); callbacksI != nil {
			//	callbacks = callbacksI.(focsCallbacks.Callbacks)
			//}
			return actions.NewTagUseCase(tagRepository, copierService, logger), nil
		},
	},
	{
		Name: "focus.page.actions.video",
		Build: func(ctn di.Container) (interface{}, error) {
			media := ctn.Get("focus.media.actions.media").(*media_usecase.Medias)
			logger := ctn.Get("logger").(*zap.SugaredLogger)
			yandexApiKey := ctn.Get("yandex.api.key").(string)
			return actions.NewVideoUseCase(media, logger, yandexApiKey), nil
		},
	},
}
