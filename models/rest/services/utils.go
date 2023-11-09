package services

import (
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetOffsetAndLimit получение офсета и лимита из query параметра
func GetOffsetAndLimit(c *gin.Context) (offset int, limit int, err error) {
	stringOffset, ok := c.GetQuery("offset")
	if !ok {
		return 0, 0, errors.BadRequest.New("offset is required")
	}
	offset, err = strconv.Atoi(stringOffset)
	if err != nil {
		return 0, 0, errors.BadRequest.Wrap(err, "error converting offset to int")
	}

	stringLimit, ok := c.GetQuery("limit")
	if !ok {
		return 0, 0, errors.BadRequest.New("offset is required")
	}
	limit, err = strconv.Atoi(stringLimit)
	if err != nil {
		return 0, 0, errors.BadRequest.Wrap(err, "error converting limit to int")
	}

	return
}
