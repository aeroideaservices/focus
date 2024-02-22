package plugin

import (
	"net/url"

	"github.com/aeroideaservices/focus/media/plugin/actions"
	"github.com/aeroideaservices/focus/media/plugin/service"
	focsCallbacks "github.com/aeroideaservices/focus/services/callbacks"
	"github.com/sarulabs/di/v2"
)

var Definitions = []di.Def{
	{
		Build: func(ctn di.Container) (interface{}, error) {
			mediaRepository := ctn.Get("focus.media.repository.media").(actions.MediaRepository)
			folderRepository := ctn.Get("focus.media.repository.folder").(actions.FolderRepository)
			mediaStorage := ctn.Get("focus.media.fileStorage").(actions.FileStorage)
			mediaProvider := ctn.Get("focus.media.provider").(*service.MediaProvider)

			var callbacks focsCallbacks.Callbacks
			if callbacksI, _ := ctn.SafeGet("focus.media.actions.media.callbacks"); callbacksI != nil {
				callbacks = callbacksI.(focsCallbacks.Callbacks)
			}

			return actions.NewMedias(mediaRepository, folderRepository, mediaStorage, mediaProvider, callbacks), nil
		},
		Name: "focus.media.actions.media",
	},
	{
		Build: func(ctn di.Container) (interface{}, error) {
			mediaRepository := ctn.Get("focus.media.repository.media").(actions.MediaRepository)
			folderRepository := ctn.Get("focus.media.repository.folder").(actions.FolderRepository)
			mediaStorage := ctn.Get("focus.media.fileStorage").(actions.FileStorage)
			mediaProvider := ctn.Get("focus.media.provider").(*service.MediaProvider)

			var callbacks focsCallbacks.Callbacks
			if callbacksI, _ := ctn.SafeGet("focus.media.actions.folders.callbacks"); callbacksI != nil {
				callbacks = callbacksI.(focsCallbacks.Callbacks)
			}

			return actions.NewFolders(folderRepository, mediaRepository, mediaStorage, mediaProvider, callbacks), nil
		},
		Name: "focus.media.actions.folder",
	},
	{
		Build: func(ctn di.Container) (interface{}, error) {
			mediaRepository := ctn.Get("focus.media.repository.media").(actions.MediaRepository)
			folderRepository := ctn.Get("focus.media.repository.folder").(actions.FolderRepository)
			mediaBaseUrl := ctn.Get("focus.media.baseUrl").(*url.URL)

			return service.NewMediaProvider(mediaRepository, folderRepository, *mediaBaseUrl, url.URL{}, url.URL{}), nil
		},
		Name: "focus.media.provider",
	},
	{
		Build: func(ctn di.Container) (interface{}, error) {
			mediaRepository := ctn.Get("focus.media.repository.media").(actions.MediaRepository)
			folderRepository := ctn.Get("focus.media.repository.folder").(actions.FolderRepository)
			mediaBaseUrl := ctn.Get("focus.media.baseUrl").(*url.URL)
			proxyMediaUrl := ctn.Get("focus.media.proxyUrl").(*url.URL)
			if proxyMediaUrl == nil {
				proxyMediaUrl = &url.URL{}
			}

			return service.NewMediaProvider(mediaRepository, folderRepository, *mediaBaseUrl, *proxyMediaUrl, url.URL{}), nil
		},
		Name: "focus.media.providerWithProxy",
	},
	{
		Build: func(ctn di.Container) (interface{}, error) {
			mediaRepository := ctn.Get("focus.media.repository.media").(actions.MediaRepository)
			folderRepository := ctn.Get("focus.media.repository.folder").(actions.FolderRepository)
			mediaBaseUrl := ctn.Get("focus.media.baseUrl").(*url.URL)
			cdnMediaUrl := ctn.Get("focus.media.cdnUrl").(*url.URL)
			if cdnMediaUrl == nil {
				cdnMediaUrl = &url.URL{}
			}

			return service.NewMediaProvider(mediaRepository, folderRepository, *mediaBaseUrl, url.URL{}, *cdnMediaUrl), nil
		},
		Name: "focus.media.providerWithCdn",
	},
}
