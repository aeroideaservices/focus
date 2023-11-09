package rest

import (
	"github.com/aeroideaservices/focus/media/rest/handlers"
	"github.com/aeroideaservices/focus/media/rest/services"
	"github.com/gin-gonic/gin"
)

// Router сервис роутинга
type Router struct {
	folderHandler *handlers.FolderHandler
	mediaHandler  *handlers.MediaHandler
	errorHandler  services.ErrorHandler
}

// NewRouter конструктор
func NewRouter(folderHandler *handlers.FolderHandler,
	mediaHandler *handlers.MediaHandler,
	errorHandler services.ErrorHandler,
) *Router {
	return &Router{
		folderHandler: folderHandler,
		mediaHandler:  mediaHandler,
		errorHandler:  errorHandler,
	}
}

// SetRoutes проставление роутов
func (r *Router) SetRoutes(group *gin.RouterGroup) {
	media := group.Group("media")
	media.Use(r.errorHandler.Handle) // отлов ошибок

	media.GET("", r.folderHandler.GetAll)

	folders := media.Group("folders")
	folders.GET("", r.folderHandler.GetTree)
	folders.POST("", r.folderHandler.Create)

	folder := folders.Group(":" + handlers.FolderIdParam)
	folder.GET("", r.folderHandler.Get)
	folder.DELETE("", r.folderHandler.Delete)
	folder.PATCH("move", r.folderHandler.Move)
	folder.PATCH("rename", r.folderHandler.Rename)

	files := media.Group("files")
	files.POST("", r.mediaHandler.Create)
	files.POST("upload", r.mediaHandler.Upload)
	files.POST("upload-list", r.mediaHandler.UploadList)

	file := files.Group(":" + handlers.FileIdParam)
	file.GET("", r.mediaHandler.Get)
	file.DELETE("", r.mediaHandler.Delete)
	file.PATCH("move", r.mediaHandler.Move)
	file.PATCH("rename", r.mediaHandler.Rename)
}
