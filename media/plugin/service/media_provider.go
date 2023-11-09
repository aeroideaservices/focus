package service

import (
	"context"
	"github.com/aeroideaservices/focus/media/plugin/actions"
	"github.com/google/uuid"
	"net/url"
)

// MediaProvider сервис работы с путями медиа файлов
type MediaProvider struct {
	mediaRepository  actions.MediaRepository
	folderRepository actions.FolderRepository
	baseMediaUrl     url.URL
	proxyMediaUrl    url.URL
}

// NewMediaProvider конструктор
func NewMediaProvider(
	mediaRepository actions.MediaRepository,
	folderRepository actions.FolderRepository,
	baseMediaUrl url.URL,
	proxyMediaUrl url.URL,
) *MediaProvider {
	return &MediaProvider{
		mediaRepository:  mediaRepository,
		folderRepository: folderRepository,
		baseMediaUrl:     baseMediaUrl,
		proxyMediaUrl:    proxyMediaUrl,
	}
}

// GetUrlByFilepath получение полного пути до файла по подпути
func (p MediaProvider) GetUrlByFilepath(mediaFilepath string) string {
	if p.proxyMediaUrl.String() != "" {
		p.proxyMediaUrl.Path += "/"
		p.proxyMediaUrl.RawQuery += "file=" + mediaFilepath
		return p.proxyMediaUrl.String()
	}

	p.baseMediaUrl.Path += "/" + mediaFilepath
	return p.baseMediaUrl.String()
}

// GetUrlById получение пути до файла по id медиа
func (p MediaProvider) GetUrlById(mediaId uuid.UUID) (string, error) {
	ctx := context.Background()
	media, err := p.mediaRepository.Get(ctx, mediaId)
	if err != nil {
		return "", err
	}

	mediaFilepath := media.Filename

	if media.FolderId != nil {
		path, err := p.folderRepository.GetFolderPath(ctx, *media.FolderId)
		if err != nil {
			return "", err
		}

		mediaFilepath = path + "/" + mediaFilepath
	}

	return p.GetUrlByFilepath(mediaFilepath), nil
}
