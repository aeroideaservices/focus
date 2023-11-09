package services

import (
	"context"
	"github.com/gin-gonic/gin"
)

// ErrorHandler Error Сервис обработки ошибок
type ErrorHandler interface {
	Handle(c *gin.Context)
}

// Validator сервис валидации
type Validator interface {
	Validate(ctx context.Context, value any) error
}
