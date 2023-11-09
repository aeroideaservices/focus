package access_control

import (
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"strings"
)

type CreateActionRule func(c *gin.Context) *Action

type AccessMiddleware struct {
	accessControl *AccessControl
	newAction     CreateActionRule
}

func NewAccessMiddleware(accessControl *AccessControl, newAction CreateActionRule) *AccessMiddleware {
	return &AccessMiddleware{
		accessControl: accessControl,
		newAction:     newAction,
	}
}

func (m AccessMiddleware) CheckAccess(c *gin.Context) {
	tokenString := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if tokenString == "" {
		abortWithErr(c, errors.Unauthorized.New("auth required"))
		return
	}

	claims, err := m.accessControl.GetClaims(tokenString)
	if err != nil {
		abortWithErr(c, err)
		return
	}

	role := NewRole("", strings.Split(claims.Scope, " "))

	action := m.newAction(c)
	if role == nil || !role.HasAccess(action) {
		abortWithErr(c, errors.Forbidden.New("forbidden"))
		return
	}

	c.Set("role", role)
	c.Set("user-full-name", strings.TrimSpace(claims.LastName+" "+claims.FirstName+" "+claims.MiddleName))
	c.Set("user-id", claims.Subject)
	c.Next()
}

func abortWithErr(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}
