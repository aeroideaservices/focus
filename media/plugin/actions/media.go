package actions

import (
	"context"
	"mime"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/google/uuid"

	"github.com/aeroideaservices/focus/media/plugin/entity"
	"github.com/aeroideaservices/focus/media/plugin/service/utils"
	"github.com/aeroideaservices/focus/services/callbacks"
	"github.com/aeroideaservices/focus/services/errors"
)

const (
	maxFileSize = 204857600
)

// Medias сервис работы с медиа
type Medias struct {
	callbacks.Callbacks
	mediaRepository  MediaRepository
	folderRepository FolderRepository
	storage          FileStorage
	mediaProvider    MediaProvider
}

// NewMedias конструктор
func NewMedias(
	mediaRepository MediaRepository,
	folderRepository FolderRepository,
	storage FileStorage,
	mediaProvider MediaProvider,
	callbacks callbacks.Callbacks,
) *Medias {
	return &Medias{
		mediaRepository:  mediaRepository,
		folderRepository: folderRepository,
		storage:          storage,
		mediaProvider:    mediaProvider,
		Callbacks:        callbacks,
	}
}

// List получение списка медиа превью
func (m Medias) List(ctx context.Context, dto ListMediasShorts) (*MediaShortList, error) {
	entities, err := m.mediaRepository.GetShortList(ctx, dto.Ids)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting medias by ids")
	}

	res := make([]MediaShort, len(entities))
	for i, media := range entities {
		res[i] = MediaShort{
			Id:    media.Id,
			Url:   m.mediaProvider.GetUrlByFilepath(media.Filepath),
			Alt:   media.Alt,
			Title: media.Title,
		}
	}

	return &MediaShortList{Items: res}, nil
}

// Create создание медиа
func (m Medias) Create(ctx context.Context, action CreateMedia) (*uuid.UUID, error) {
	if action.Size > maxFileSize {
		return nil, ErrMaxFileSize
	}

	var folderPath string
	var err error
	if action.FolderId != nil {
		hasFolder := m.folderRepository.Has(ctx, *action.FolderId)
		if !hasFolder {
			return nil, ErrFolderNotFound
		}

		folderPath, err = m.folderRepository.GetFolderPath(ctx, *action.FolderId)
		if err != nil {
			return nil, errors.NoType.Wrap(err, "error getting folder path")
		}
	}

	hasMedia, mediaId := m.mediaRepository.HasByFilterWithId(
		ctx, MediaFilter{
			FolderId:     action.FolderId,
			WithFolderId: true,
			Filename:     action.Filename,
		},
	)
	if hasMedia {
		mediaFilepath := filepath.Join(folderPath, action.Filename)
		saveMediaFile := &UploadFile{
			Key:         mediaFilepath,
			ContentType: mime.TypeByExtension(filepath.Ext(action.Filename)),
			File:        action.File,
		}
		err = m.storage.Upload(ctx, saveMediaFile)
		if err != nil {
			return nil, errors.NoType.Wrap(err, "error uploading media file")
		}

		// updateMediaDto := &UpdateMediaDto{
		// 	Id:       mediaId,
		// 	Name:     strings.TrimSuffix(action.Filename, filepath.Ext(action.Filename)),
		// 	Filename: action.Filename,
		// 	Filepath: mediaFilepath,
		// 	FolderId: action.FolderId,
		// }
		// err = m.mediaRepository.Update(ctx, updateMediaDto)
		// if err != nil {
		// 	return errors.NoType.Wrap(err, "error updating media")
		// }

		return &mediaId, nil
		// return nil, ErrMediaAlreadyExistsInFolder
	}

	mediaFilepath := filepath.Join(folderPath, action.Filename)
	saveMediaFile := &UploadFile{
		Key:         mediaFilepath,
		ContentType: mime.TypeByExtension(filepath.Ext(action.Filename)),
		File:        action.File,
	}
	err = m.storage.Upload(ctx, saveMediaFile)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error uploading media file")
	}

	newId := uuid.New()
	media := entity.Media{
		Id:       newId,
		Name:     strings.TrimSuffix(action.Filename, filepath.Ext(action.Filename)),
		Filename: action.Filename,
		Alt:      action.Alt,
		Title:    action.Title,
		Size:     action.Size,
		Filepath: mediaFilepath,
		FolderId: action.FolderId,
	}
	err = m.mediaRepository.Create(ctx, media)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error creating media")
	}

	m.GoAfterCreate(media.Id)

	return &newId, nil
}

// Upload загрузка нового медиа
func (m Medias) Upload(ctx context.Context, dto CreateMedia) (string, error) {
	id, err := m.Create(ctx, dto)
	if err != nil {
		return "", err
	}

	return m.mediaProvider.GetUrlById(*id)
}

// UploadList загрузка нескольких медиа
func (m Medias) UploadList(ctx context.Context, dto CreateMediasList) ([]uuid.UUID, error) {
	for _, file := range dto.Files {
		if file.Size > maxFileSize {
			return nil, ErrMaxFileSize
		}
	}

	var folderPath string
	var err error
	if dto.FolderId != nil {
		hasFolder := m.folderRepository.Has(ctx, *dto.FolderId)
		if !hasFolder {
			return nil, ErrFolderNotFound
		}

		folderPath, err = m.folderRepository.GetFolderPath(ctx, *dto.FolderId)
		if err != nil {
			return nil, errors.NoType.Wrap(err, "error getting folder path")
		}
	}

	filenames := make([]string, len(dto.Files))
	for i, file := range dto.Files {
		filenames[i] = file.Filename
	}

	//hasMedia := m.mediaRepository.HasByFilter(ctx, MediaFilter{
	//	FolderId:     dto.FolderId,
	//	WithFolderId: true,
	//	Filenames:    filenames,
	//})
	//if hasMedia {
	//	return nil, ErrMediaAlreadyExistsInFolder
	//}

	createMediaFiles := make([]UploadFile, len(dto.Files))
	for i, file := range dto.Files {
		createMediaFiles[i] = UploadFile{
			Key:         filepath.Join(folderPath, file.Filename),
			ContentType: mime.TypeByExtension(filepath.Ext(file.Filename)),
			File:        file.File,
		}
	}

	err = m.storage.UploadList(ctx, createMediaFiles...)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error uploading media file")
	}

	entities := make([]entity.Media, 0)
	ids := make([]uuid.UUID, len(dto.Files))
	for i, f := range dto.Files {
		hasMedia, mediaId := m.mediaRepository.HasByFilterWithId(
			ctx, MediaFilter{
				FolderId:     dto.FolderId,
				WithFolderId: true,
				Filename:     f.Filename,
			},
		)
		if hasMedia {
			ids[i] = mediaId
		} else {
			newId := uuid.New()
			ids[i] = newId
			entities = append(
				entities, entity.Media{
					Id:       newId,
					Name:     strings.TrimSuffix(f.Filename, filepath.Ext(f.Filename)),
					Filename: f.Filename,
					Size:     f.Size,
					Filepath: createMediaFiles[i].Key,
					FolderId: dto.FolderId,
				},
			)
		}
	}

	if len(entities) != 0 {
		err = m.mediaRepository.Create(ctx, entities...)
		if err != nil {
			return nil, errors.NoType.Wrap(err, "error creating medias")
		}
	}

	m.GoAfterCreate(ids...)

	return ids, nil
}

// Get получение медиа по id
func (m Medias) Get(ctx context.Context, dto GetMedia) (*MediaPreview, error) {
	media, err := m.mediaRepository.Get(ctx, dto.Id)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting media by id")
	}

	res := &MediaPreview{
		Id:          media.Id,
		Name:        media.Name,
		Ext:         strings.TrimPrefix(filepath.Ext(media.Filename), "."),
		Size:        utils.Filesize(media.Size),
		Alt:         media.Alt,
		Title:       media.Title,
		ContentType: mime.TypeByExtension(filepath.Ext(media.Filename)),
		Url:         m.mediaProvider.GetUrlByFilepath(media.Filepath),
		UpdatedAt:   utils.Time(media.UpdatedAt),
		FolderId:    media.FolderId,
	}

	return res, nil
}

// Rename переименование медиа
func (m Medias) Rename(ctx context.Context, dto RenameMedia) error {
	media, err := m.mediaRepository.Get(ctx, dto.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting media by id")
	}
	if media.Name == dto.Name {
		return ErrMediaAlreadyHasSameName
	}

	folderPath := ""
	if media.FolderId != nil {
		folderPath, err = m.folderRepository.GetFolderPath(ctx, *media.FolderId)
		if err != nil {
			return errors.NoType.Wrap(err, "error getting folder path")
		}
	}

	newFilename := dto.Name + filepath.Ext(media.Filename)
	hasMedia := m.mediaRepository.HasByFilter(
		ctx, MediaFilter{
			Filename:     newFilename,
			FolderId:     media.FolderId,
			WithFolderId: true,
		},
	)
	if hasMedia {
		return ErrMediaAlreadyExistsInFolder
	}

	newFilepath := filepath.Join(folderPath, newFilename)
	oldFilepath := media.Filepath
	err = m.storage.Move(ctx, oldFilepath, newFilepath)
	if err != nil {
		return errors.NoType.Wrap(err, "error moving media file")
	}

	updateMediaDto := &UpdateMediaDto{
		Id:       media.Id,
		Name:     dto.Name,
		Filename: newFilename,
		Filepath: newFilepath,
		FolderId: media.FolderId,
	}
	err = m.mediaRepository.Update(ctx, updateMediaDto)
	if err != nil {
		return errors.NoType.Wrap(err, "error updating media")
	}

	m.GoAfterUpdate(media.Id)

	return nil
}

// Move перемещение медиа
func (m Medias) Move(ctx context.Context, dto MoveMedia) error {
	media, err := m.mediaRepository.Get(ctx, dto.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting media by id")
	}
	if reflect.DeepEqual(dto.FolderId, media.FolderId) {
		return ErrMediaAlreadyHasSameFolder
	}

	folderPath := ""
	if dto.FolderId != nil {
		hasFolder := m.folderRepository.Has(ctx, *dto.FolderId)
		if !hasFolder {
			return ErrFolderNotFound
		}

		folderPath, err = m.folderRepository.GetFolderPath(ctx, *dto.FolderId)
		if err != nil {
			return errors.NoType.Wrap(err, "error getting folder path")
		}
	}

	hasMedia := m.mediaRepository.HasByFilter(
		ctx, MediaFilter{
			Filename:     media.Filename,
			FolderId:     dto.FolderId,
			WithFolderId: true,
		},
	)
	if hasMedia {
		return ErrMediaAlreadyExistsInFolder
	}

	oldFilepath := media.Filepath
	newFilepath := filepath.Join(folderPath, media.Filename)
	err = m.storage.Move(ctx, oldFilepath, newFilepath)
	if err != nil {
		return errors.NoType.Wrap(err, "error moving media file")
	}

	updateMediaDto := &UpdateMediaDto{
		Id:       media.Id,
		Name:     media.Name,
		Filename: media.Filename,
		Filepath: newFilepath,
		FolderId: dto.FolderId,
	}
	err = m.mediaRepository.Update(ctx, updateMediaDto)
	if err != nil {
		return errors.NoType.Wrap(err, "error updating media")
	}

	m.GoAfterUpdate(media.Id)

	return nil
}

// Delete удаление медиа по id
func (m Medias) Delete(ctx context.Context, dto GetMedia) error {
	media, err := m.mediaRepository.Get(ctx, dto.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting media by id")
	}

	err = m.storage.Delete(ctx, media.Filepath)
	if err != nil {
		return errors.NoType.Wrap(err, "error deleting media file")
	}

	err = m.mediaRepository.Delete(ctx, dto.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error deleting media")
	}

	m.GoAfterDelete(media.Id)

	return nil
}

// CheckIds Проверяет существование медиа с такими id
func (m Medias) CheckIds(ctx context.Context, ids ...uuid.UUID) error {
	count, err := m.mediaRepository.Count(ctx, MediaFilter{InIds: ids})
	if err != nil {
		return errors.NoType.Wrap(err, "error counting medias")
	}
	if count != len(ids) {
		return ErrOneOfMediasNotExists
	}

	return nil
}

func (m Medias) Download(ctx context.Context, dto GetMedia) (string, error) {
	media, err := m.mediaRepository.Get(ctx, dto.Id)
	if err != nil {
		return "", errors.NoType.Wrap(err, "error getting media by id")
	}
	err = m.storage.DownloadFile(ctx, media.Filepath, media.Filename)
	if err != nil {
		return "", errors.NoType.Wrap(err, "error downloading media file")
	}

	return media.Filename, nil
}

func (m Medias) UpdateSubtitles(ctx context.Context, dto UpdateMediaSubtitles) error {

	media, err := m.mediaRepository.Get(ctx, dto.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting media by id")
	}

	err = m.mediaRepository.UpdateSubtitles(ctx, dto.Id, dto.Subtitles)
	if err != nil {
		return errors.NoType.Wrap(err, "error updating media")
	}

	m.GoAfterUpdate(media.Id)

	return nil
}
