package gin

import (
	"net/http"

	"github.com/aeroideaservices/focus/services/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	errorStatuses = map[errors.ErrorType]int{
		errors.NoType:             http.StatusInternalServerError,
		errors.BadRequest:         http.StatusBadRequest,
		errors.NotFound:           http.StatusNotFound,
		errors.Conflict:           http.StatusConflict,
		errors.RequestTimeout:     http.StatusRequestTimeout,
		errors.ServiceUnavailable: http.StatusServiceUnavailable,
	}

	errorMessages = map[errors.ErrorType]string{
		errors.NoType:             "Произошла ошибка. Попробуйте выполнить операцию позже.",
		errors.BadRequest:         "Некорректный запрос, отсутствует один из обязательных параметров или один из параметров некорректный.",
		errors.NotFound:           "Запись не найдена.",
		errors.Conflict:           "Конфликт.",
		errors.RequestTimeout:     "Превышено время ожидания запроса.",
		errors.ServiceUnavailable: "В настоящий момент сервис недоступен.",
	}
)

type Translatable interface {
	Translate(translator *message.Printer) string
	Error() string
}

type ErrorHandler struct {
	logger *zap.SugaredLogger
}

func NewErrorHandler(defaultLanguage string, logger *zap.SugaredLogger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
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

	h.log(err)

	responseData, status := getResponseError(getLanguage(c), err)
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

func getLanguage(c *gin.Context) language.Tag {
	lang, _ := c.Cookie("lang")
	accept := c.GetHeader("accept-language")
	fallback := "ru"
	return message.MatchLanguage(lang, accept, fallback)
}

func getResponseError(lang language.Tag, err error) (errorFormatted responseError, statusCode int) {
	errorType := errors.GetType(err)
	if errStatus, ok := errorStatuses[errorType]; ok {
		statusCode = errStatus
	} else {
		statusCode = http.StatusInternalServerError
	}

	statusText := http.StatusText(statusCode)

	var msg string
	if e, ok := err.(Translatable); ok {
		msg = e.Translate(message.NewPrinter(lang))
	}
	if msg == "" {
		msg = errorMessages[errorType]
	}

	errorFormatted.Error = statusText
	errorFormatted.Code = statusCode
	errorFormatted.Message = msg
	errorFormatted.Debug = err.Error()

	return errorFormatted, statusCode
}
