package actions

import "github.com/aeroideaservices/focus/services/errors"

var (
	errMediaPluginIsNotImported = errors.Internal.New("media plugin is not imported").T("media-plugin-not-imported")
	errModelElementConflict     = errors.Conflict.Newf("model element with the same value of field already exists")
)
