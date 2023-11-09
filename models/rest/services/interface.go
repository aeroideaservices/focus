package services

import (
	"context"
	"github.com/gin-gonic/gin"
)

// ErrorHandler Error Сервис обработки ошибок
type ErrorHandler interface {
	Handle(c *gin.Context)
}

// ErrorService Error Сервис обработки ошибок
type ErrorService interface {
	HandleError(c *gin.Context, err error)
}

type Validator interface {
	Validate(ctx context.Context, value any) error
}
