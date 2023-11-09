package actions

import (
	"github.com/aeroideaservices/focus/services/errors"
)

var (
	ErrFolderAlreadyExists        = errors.Conflict.New("folder already exists").T("folder.conflict")
	ErrFolderRecursiveAttachment  = errors.Conflict.New("folders recursion attaching").T("folder.recursion")
	ErrFolderAlreadyInThisFolder  = errors.BadRequest.New("folder already in this folder").T("folder.moving-same-folder")
	ErrFolderAlreadyHasSameName   = errors.BadRequest.New("folder already has the same name").T("folder.renaming-same-name")
	ErrMediaAlreadyHasSameName    = errors.BadRequest.New("media already has the same name").T("media.renaming-same-name")
	ErrMediaAlreadyHasSameFolder  = errors.BadRequest.New("media already in this folder").T("media.moving-same-folder")
	ErrMediaAlreadyExistsInFolder = errors.Conflict.New("folder with the same name already exists in this folder").T("media.conflict")
	ErrFolderNotFound             = errors.NotFound.New("folder not found").T("folder.not-found")
	ErrMediaNotFound              = errors.NotFound.New("media not found").T("media.not-found")
	ErrOneOfMediasNotExists       = errors.BadRequest.New("some medias do not exist").T("media.not-exists")
	ErrMaxFileSize                = errors.BadRequest.New("media file size too large").T("media.file.size")
)
