package services

import (
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetLimitAndOffset получение лимита и офсета из контекста
func GetLimitAndOffset(ctx *gin.Context, limit *int, offset *int) error {
	var err error

	stringLimit, hasLimit := ctx.GetQuery("limit")
	if !hasLimit {
		return errors.BadRequest.New("limit must be non-empty")
	}
	*limit, err = strconv.Atoi(stringLimit)
	if err != nil {
		return errors.BadRequest.Wrap(err, "limit must be integer")
	}

	stringOffset, hasOffset := ctx.GetQuery("offset")
	if !hasOffset {
		return errors.BadRequest.New("offset must be non-empty")
	}
	*offset, err = strconv.Atoi(stringOffset)
	if err != nil {
		return errors.BadRequest.Wrap(err, "offset must be integer")
	}

	return nil
}
