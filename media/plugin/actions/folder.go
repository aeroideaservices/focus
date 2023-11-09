package actions

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/aeroideaservices/focus/media/plugin/entity"
	"github.com/aeroideaservices/focus/services/callbacks"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/google/uuid"
)

// Folders сервис работы с папками
type Folders struct {
	callbacks.Callbacks
	folderRepository FolderRepository
	mediaRepository  MediaRepository
	fileStorage      FileStorage
	mediaProvider    MediaProvider
}

// NewFolders конструктор
func NewFolders(
	folderRepository FolderRepository,
	mediaRepository MediaRepository,
	mediaStorage FileStorage,
	mediaProvider MediaProvider,
	callbacks callbacks.Callbacks,
) *Folders {
	return &Folders{
		folderRepository: folderRepository,
		mediaRepository:  mediaRepository,
		fileStorage:      mediaStorage,
		mediaProvider:    mediaProvider,
		Callbacks:        callbacks,
	}
}

// GetAll получение папок и медиа, лежащих в одной папке
func (f Folders) GetAll(ctx context.Context, filter FolderFilter) (*FolderAndMediasList, error) {
	if filter.Filter.FolderId != nil {
		hasFolder := f.folderRepository.Has(ctx, *filter.Filter.FolderId)
		if !hasFolder {
			return nil, ErrFolderNotFound
		}
	}

	list, err := f.folderRepository.GetFoldersAndMedias(ctx, filter)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting folders and medias")
	}

	res := &FolderAndMediasList{
		Total: list.Total,
		Items: make([]FolderAndMedias, len(list.Items)),
	}
	for i, item := range list.Items {
		var fam = FolderAndMedias{ResourceType: item.ResourceType}
		switch item.ResourceType {
		case "folder":
			fam.FolderFields = &FolderFields{
				Id:   item.Id,
				Name: item.Name,
				Size: item.Size,
			}
		case "file":
			fam.FileFields = &FileFields{
				Id:   item.Id,
				Name: item.Name,
				Size: item.Size,
				Url:  f.mediaProvider.GetUrlByFilepath(item.Filepath),
				Ext:  strings.TrimPrefix(filepath.Ext(item.Filepath), "."),
			}
		default:
			return nil, errors.BadRequest.Newf("wrong resource type, got %v, expected file or folder", item.ResourceType)
		}
		res.Items[i] = fam
	}

	folders, err := f.folderRepository.GetFolderParents(ctx, filter.Filter)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting folder parents")
	}

	res.Breadcrumbs = make([]FolderBreadcrumb, 0)
	for _, folder := range folders {
		res.Breadcrumbs = append(res.Breadcrumbs, FolderBreadcrumb{
			Name:     folder.Name,
			FolderId: pointer(folder.Id),
		})
	}

	return res, nil
}

func pointer[T any](val T) *T {
	return &val
}

// GetTree получение дерева папок
func (f Folders) GetTree(ctx context.Context) ([]*FolderResponse, error) {
	folders, err := f.folderRepository.GetFoldersTree(ctx)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting folders tree")
	}

	return folders, nil
}

// Get получение папки
func (f Folders) Get(ctx context.Context, action GetFolder) (*FolderDetail, error) {
	folder, err := f.folderRepository.GetWithSize(ctx, action.Id)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error getting folder by id")
	}

	return folder, nil
}

// Create создание папки
func (f Folders) Create(ctx context.Context, action CreateFolder) (*uuid.UUID, error) {
	if action.ParentFolderId != nil {
		hasFolder := f.folderRepository.Has(ctx, *action.ParentFolderId)
		if !hasFolder {
			return nil, ErrFolderNotFound
		}
	}

	hasFolder := f.folderRepository.HasByFilter(ctx, Filter{
		Name:         action.Name,
		FolderId:     action.ParentFolderId,
		WithFolderId: true,
	})
	if hasFolder {
		return nil, ErrFolderAlreadyExists
	}

	newId := uuid.New()
	folder := &entity.Folder{
		Id:       newId,
		Name:     action.Name,
		FolderId: action.ParentFolderId,
	}
	err := f.folderRepository.Create(ctx, folder)
	if err != nil {
		return nil, errors.NoType.Wrap(err, "error creating folder")
	}

	f.GoAfterCreate(folder.Id)

	return &newId, nil
}

// Rename переименование папки
func (f Folders) Rename(ctx context.Context, action RenameFolder) error {
	folder, err := f.folderRepository.Get(ctx, action.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting folder by id")
	}
	if folder.Name == action.Name {
		return ErrFolderAlreadyHasSameName
	}

	hasFolder := f.folderRepository.HasByFilter(ctx, Filter{
		Name:         action.Name,
		FolderId:     folder.FolderId,
		WithFolderId: true,
	})
	if hasFolder {
		return ErrFolderAlreadyExists
	}

	folder.Name = action.Name
	err = f.folderRepository.Update(ctx, folder)
	if err != nil {
		return errors.NoType.Wrap(err, "error updating folder name")
	}

	medias, err := f.folderRepository.GetAllFolderMedias(ctx, action.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting folder medias")
	}

	if len(medias) != 0 {
		for _, media := range medias {
			err := f.fileStorage.Move(ctx, media.Filepath, media.NewFilepath)
			if err != nil {
				return errors.NoType.Wrap(err, "error moving folder media file")
			}
			media.Filepath = media.NewFilepath
		}
		if err = f.mediaRepository.Update(ctx, medias...); err != nil {
			return errors.NoType.Wrap(err, "error updating medias")
		}
	}

	f.GoAfterUpdate(folder.Id)

	return nil
}

// Move перемещение папки
func (f Folders) Move(ctx context.Context, action MoveFolder) error {
	folder, err := f.folderRepository.Get(ctx, action.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting folder by id")
	}
	if folder.FolderId == action.ParentFolderId {
		return ErrFolderAlreadyInThisFolder
	}

	if action.ParentFolderId != nil {
		hasFolder := f.folderRepository.Has(ctx, *action.ParentFolderId)
		if !hasFolder {
			return ErrFolderNotFound
		}
	}

	hasSubFolder, err := f.folderRepository.HasSubFolder(ctx, action.Id, action.ParentFolderId)
	if hasSubFolder {
		return ErrFolderRecursiveAttachment
	}

	hasFolder := f.folderRepository.HasByFilter(ctx, Filter{
		Name:         folder.Name,
		FolderId:     action.ParentFolderId,
		WithFolderId: true,
	})
	if hasFolder {
		return ErrFolderAlreadyExists
	}

	folder.FolderId = action.ParentFolderId
	err = f.folderRepository.Update(ctx, folder)
	if err != nil {
		return errors.NoType.Wrap(err, "error updating folder parent id")
	}

	medias, err := f.folderRepository.GetAllFolderMedias(ctx, action.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting folder medias")
	}

	if len(medias) != 0 {
		for _, media := range medias {
			err = f.fileStorage.Move(ctx, media.Filepath, media.NewFilepath)
			if err != nil {
				return errors.NoType.Wrap(err, "error moving folder media file")
			}

			media.Filepath = media.NewFilepath
		}
		err = f.mediaRepository.Update(ctx, medias...)
		if err != nil {
			return errors.NoType.Wrap(err, "error updating folder medias")
		}
	}

	f.GoAfterUpdate(folder.Id)

	return nil
}

// Delete удаление папко по id
func (f Folders) Delete(ctx context.Context, action GetFolder) error {
	hasFolder := f.folderRepository.Has(ctx, action.Id)
	if !hasFolder {
		return ErrFolderNotFound
	}

	filePaths, err := f.folderRepository.GetFolderMediaFilePaths(ctx, &action.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error getting media file path")
	}
	if len(filePaths) > 0 {
		err = f.fileStorage.Delete(ctx, filePaths...)
		if err != nil {
			return errors.NoType.Wrap(err, "error remove file")
		}
	}

	err = f.folderRepository.Delete(ctx, action.Id)
	if err != nil {
		return errors.NoType.Wrap(err, "error remove folder")
	}

	f.GoAfterDelete(action.Id)

	return nil
}
