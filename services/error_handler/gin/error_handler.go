package gin

import (
	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
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

type ErrorHandler struct {
	logger *zap.SugaredLogger
}

func NewErrorHandler(logger *zap.SugaredLogger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}

type responseError struct {
	Code    int         `json:"code"`
	Error   interface{} `json:"error"`
	Message string      `json:"message"`
	Debug   interface{} `json:"debug"`
}

func (h ErrorHandler) HandleError(c *gin.Context, err error) {
	var status int
	errorType := errors.GetType(err)
	if errStatus, ok := errorStatuses[errorType]; ok {
		status = errStatus
	} else {
		status = http.StatusInternalServerError
	}

	h.log(err)

	responseData := getResponseError(err)
	c.JSON(status, responseData)

}

func (h ErrorHandler) log(err error) {
	errorType := errors.GetType(err)
	stackTrace := errors.GetStackTrace(errors.Cause(err))

	var logFunc func(msg string, keysAndValues ...interface{})
	switch errorType {
	case errors.NoType:
		logFunc = h.logger.Errorw
	default:
		logFunc = h.logger.Debugw
	}

	logFunc("service error", "error", err.Error(), "stackTrace", stackTrace)
}

func getResponseError(err error) responseError {
	errorFormatted := responseError{}

	status := errorStatuses[errors.GetType(err)]
	statusText := http.StatusText(status)
	message := errorMessages[errors.GetType(err)]

	errorFormatted.Error = statusText
	errorFormatted.Code = status
	errorFormatted.Message = message
	errorFormatted.Debug = err.Error()

	return errorFormatted
}
