package middleware

import (
	"net/http"

	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
)

var (
	errorStatuses = map[errors.ErrorType]int{
		errors.NoType:             http.StatusInternalServerError,
		errors.BadRequest:         http.StatusBadRequest,
		errors.Unauthorized:       http.StatusUnauthorized,
		errors.Forbidden:          http.StatusForbidden,
		errors.NotFound:           http.StatusNotFound,
		errors.RequestTimeout:     http.StatusRequestTimeout,
		errors.Conflict:           http.StatusConflict,
		errors.ServiceUnavailable: http.StatusServiceUnavailable,
	}

	errorMessages = map[errors.ErrorType]string{
		errors.NoType:             "Произошла ошибка. Попробуйте выполнить операцию позже.",
		errors.BadRequest:         "Некорректный запрос, отсутствует один из обязательных параметров или один из параметров некорректный.",
		errors.Unauthorized:       "Пользователь не авторизован.",
		errors.Forbidden:          "Доступ запрещен.",
		errors.NotFound:           "Запись не найдена.",
		errors.RequestTimeout:     "Превышено время ожидания запроса.",
		errors.Conflict:           "Конфликт.",
		errors.ServiceUnavailable: "В настоящий момент сервис недоступен.",
	}
)

type errorTranslator interface {
	Translate(err error) (string, error)
}

type ErrorHandler struct {
	logger logger
	et     errorTranslator
}

func NewErrorHandler(
	logger logger,
) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

func (e *ErrorHandler) SetTranslator(trans errorTranslator) *ErrorHandler {
	e.et = trans
	return e
}

type logger interface {
	Errorw(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Debugw(msg string, keysAndValues ...any)
}

type responseError struct {
	Code    int         `json:"code"`
	Error   interface{} `json:"error"`
	Message string      `json:"message"`
	Debug   interface{} `json:"debug"`
}

func (h ErrorHandler) Handle(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return
	}
	err := c.Errors[0].Err

	h.logErr(err)

	responseData, status := h.getResponseError(err)
	c.JSON(status, responseData)
	c.Abort()
}

func (h ErrorHandler) log(err error) {
	errorType := errors.GetType(err)
	stackTrace := errors.GetStackTrace(errors.Cause(err))

	var logFunc func(msg string, keysAndValues ...interface{})
	switch errorType {
	case errors.NoType:
		logFunc = h.logger.Errorw
	case errors.ServiceUnavailable:
		logFunc = h.logger.Warnw
	case errors.RequestTimeout:
		logFunc = h.logger.Infow
	default:
		logFunc = h.logger.Debugw
	}

	logFunc("service error", "error", err.Error(), "stackTrace", stackTrace)
}

func (h ErrorHandler) logErr(err error) {
	if errors.GetType(err) != errors.NoType {
		return
	}

	stackTrace := errors.GetStackTrace(errors.Cause(err))
	h.logger.Errorw("service error", "error", err.Error(), "stackTrace", stackTrace)
}

func (h ErrorHandler) getResponseError(err error) (errorFormatted responseError, statusCode int) {
	errorType := errors.GetType(err)
	if errStatus, ok := errorStatuses[errorType]; ok {
		statusCode = errStatus
	} else {
		statusCode = http.StatusInternalServerError
	}

	msg := errorMessages[errorType]
	if h.et != nil {
		if translation, err := h.et.Translate(err); err == nil {
			msg = translation
		}
	}

	errorFormatted.Error = http.StatusText(statusCode)
	errorFormatted.Code = statusCode
	errorFormatted.Message = msg
	errorFormatted.Debug = err.Error()

	return errorFormatted, statusCode
}
