package handlers

import (
	"github.com/aeroideaservices/focus/models/plugin/actions"
	"github.com/aeroideaservices/focus/models/rest/services"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ExportHandler обработчик запросов, связанных с экспортом элементов модели
type ExportHandler struct {
	export    *actions.Export
	validator services.Validator
}

// NewExportHandler конструктор
func NewExportHandler(
	export *actions.Export,
	validator services.Validator,
) *ExportHandler {
	return &ExportHandler{
		export:    export,
		validator: validator,
	}
}

// Export экспорт элементов модели
func (h ExportHandler) Export(c *gin.Context) {
	action := actions.ExportModelElements{}
	action.ModelCode = c.Param(ModelCodeParam)

	err := h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	err = h.export.Export(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetExportInfo получение информации по последнему экспорту
func (h ExportHandler) GetExportInfo(c *gin.Context) {
	action := actions.GetExportInfo{}
	action.ModelCode = c.Param(ModelCodeParam)

	err := h.validator.Validate(c, action)
	if err != nil {
		_ = c.Error(errors.BadRequest.Wrap(err, "error validating"))
		return
	}

	exportInfo, err := h.export.GetExportInfo(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, exportInfo)
}
