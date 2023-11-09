package actions

import (
	focus_errors "github.com/aeroideaservices/focus/services/errors"
)

var (
	ErrConfAlreadyExists        = focus_errors.Conflict.New("configuration with the same code  already exists").T("conf.exists")
	ErrOptAlreadyExists         = focus_errors.Conflict.New("option with the same code already exists").T("opt.exists")
	ErrOptLinkedToAnotherConf   = focus_errors.BadRequest.New("option linked to another configuration").T("opt.linked-to-another")
	ErrFieldNotUpdatable        = focus_errors.BadRequest.New("field is not updatable").T("field-not-updatable")
	ErrMediaPluginIsNotImported = focus_errors.BadRequest.New("media plugin is not imported").T("media-not-imported")
	ErrConfNotFound             = focus_errors.NotFound.New("configuration not found").T("conf.not-found")
	ErrOptNotFound              = focus_errors.NotFound.New("option not found").T("opt.not-found")
)
