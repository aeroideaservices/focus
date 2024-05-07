package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/jemzee04/focus/page/rest/handlers"
	"github.com/jemzee04/focus/page/rest/services"
)

type Router struct {
	pageHandler    *handlers.PageHandler
	galleryHandler *handlers.GalleryHandler
	cardHandler    *handlers.CardHandler
	tagHandler     *handlers.TagHandler
	videoHandler   *handlers.VideoHandler
	errorHandler   services.ErrorHandler
}

func NewRouter(
	pageHandler *handlers.PageHandler, galleryHandler *handlers.GalleryHandler,
	cardHandler *handlers.CardHandler, tagHandler *handlers.TagHandler,
	videoHandler *handlers.VideoHandler,
	errorHandler services.ErrorHandler,
) *Router {
	return &Router{
		pageHandler:    pageHandler,
		galleryHandler: galleryHandler,
		cardHandler:    cardHandler,
		tagHandler:     tagHandler,
		videoHandler:   videoHandler,
		errorHandler:   errorHandler,
	}
}

func (r *Router) SetRoutes(routerGroup *gin.RouterGroup) {
	pages := routerGroup.Group("pages")
	pages.Use(r.errorHandler.Handle)
	pages.POST("", r.pageHandler.Create)
	pages.GET("", r.pageHandler.GetList)
	pages.GET("/:page-id", r.pageHandler.GetById)
	pages.DELETE("/:page-id", r.pageHandler.Delete)
	pages.PATCH("/:page-id/galleries/:gallery-id", r.pageHandler.PatchGalleryPosition)
	pages.PATCH("/:page-id/galleries/link", r.pageHandler.LinkGalleries)
	pages.PATCH("/:page-id/galleries/unlink", r.pageHandler.UnlinkGalleries)
	pages.PATCH("/:page-id/properties", r.pageHandler.PatchProperties)

	galleries := pages.Group("galleries")

	galleries.GET("", r.galleryHandler.GetList)
	galleries.POST("", r.galleryHandler.Create)
	galleries.DELETE("", r.galleryHandler.DeleteList)
	galleries.GET("/:gallery-id", r.galleryHandler.GetById)
	galleries.PUT("/:gallery-id", r.galleryHandler.Update)
	galleries.PATCH("/:gallery-id/name", r.galleryHandler.PatchName)
	galleries.PATCH("/:gallery-id/card/:card-id", r.galleryHandler.PatchCardPosition)
	galleries.PATCH("/:gallery-id/card/link", r.galleryHandler.LinkCards)
	galleries.PATCH("/:gallery-id/card/unlink", r.galleryHandler.UnlinkCards)

	cards := pages.Group("cards")

	cards.GET("", r.cardHandler.GetList)
	cards.POST("", r.cardHandler.Create)
	cards.GET("/:card-id", r.cardHandler.GetById)
	cards.PUT("/:card-id", r.cardHandler.Update)
	cards.DELETE("/:card-id", r.cardHandler.Delete)
	cards.PATCH("/:card-id/user", r.cardHandler.PatchUser)
	cards.PATCH("/:card-id/previewtext", r.cardHandler.PatchPreviewText)
	cards.PATCH("/:card-id/detailtext", r.cardHandler.PatchDetailText)
	cards.PATCH("/:card-id/learn-more-url", r.cardHandler.PatchLearnMoreUrl)
	cards.PATCH("/:card-id/tags", r.cardHandler.PatchTags)
	cards.PATCH("/:card-id/tags/link", r.cardHandler.LinkTags)
	cards.PATCH("/:card-id/tags/unlink", r.cardHandler.UnlinkTags)

	tags := pages.Group("tags")

	tags.GET("", r.tagHandler.GetList)
	tags.POST("", r.tagHandler.Create)
	tags.PUT("/:tag-id", r.tagHandler.Update)
	tags.DELETE("/:tag-id", r.tagHandler.Delete)

	pages.POST("/video/upload", r.videoHandler.Create)
}
