package handlers

import (
	"github.com/aeroideaservices/focus/menu/plugin/actions"
	"github.com/aeroideaservices/focus/menu/rest/services"
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DomainsHandler struct {
	domains   *actions.Domains
	validator services.Validator
}

func NewDomainsHandler(domains *actions.Domains, validator services.Validator) *DomainsHandler {
	return &DomainsHandler{domains: domains, validator: validator}
}

func (h DomainsHandler) List(c *gin.Context) {
	action := actions.ListDomains{}
	if err := services.GetLimitAndOffset(c, &action.Limit, &action.Offset); err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.validator.Validate(c, action); err != nil {
		_ = c.Error(err)
		return
	}

	list, err := h.domains.List(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h DomainsHandler) Create(c *gin.Context) {
	action := actions.CreateDomain{}
	if err := c.ShouldBindJSON(&action); err != nil {
		_ = c.Error(errors.BadRequest.Wrapf(err, "json binding error"))
		return
	}

	if err := h.validator.Validate(c, action); err != nil {
		_ = c.Error(err)
		return
	}

	id, err := h.domains.Create(c, action)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}
