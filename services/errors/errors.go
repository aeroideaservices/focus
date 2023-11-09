package errors

import (
	"github.com/pkg/errors"
)

const (
	NoType = ErrorType(iota)
	BadRequest
	NotFound
	Conflict
	RequestTimeout
	ServiceUnavailable
	Forbidden
	Unauthorized
	//add any type you want

	Internal = NoType
)

type ErrorType uint

type FocusError struct {
	errorType     ErrorType
	originalError error
	Trans         *Trans
}

type Trans struct {
	Msg    string
	Params []string
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Error returns the message of a FocusError
func (e FocusError) Error() string {
	return e.originalError.Error()
}

// New creates a new FocusError
func (t ErrorType) New(msg string) FocusError {
	return FocusError{errorType: t, originalError: errors.New(msg)}
}

// Newf creates a new FocusError with formatted message
func (t ErrorType) Newf(msg string, args ...interface{}) FocusError {
	return FocusError{errorType: t, originalError: errors.Errorf(msg, args...)}
}

// Wrap creates a new wrapped error
func (t ErrorType) Wrap(err error, msg string) FocusError {
	return t.Wrapf(err, msg)
}

// Wrapf creates a new wrapped error with formatted message
func (t ErrorType) Wrapf(err error, msg string, args ...interface{}) FocusError {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(FocusError); ok {
		return FocusError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			Trans:         customErr.Trans,
		}
	}

	return FocusError{errorType: t, originalError: wrappedError}
}

func (e FocusError) Unwrap() error {
	return errors.Unwrap(e.originalError)
}

func (e FocusError) T(msg string, params ...string) FocusError {
	return FocusError{
		errorType:     e.errorType,
		originalError: e.originalError,
		Trans: &Trans{
			Msg:    msg,
			Params: params,
		},
	}
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(FocusError); ok {
		return customErr.errorType
	}

	return NoType
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func GetStackTrace(err error) errors.StackTrace {
	if customErr, ok := err.(FocusError); ok {
		err = customErr.originalError
	}
	if stacked, ok := err.(stackTracer); ok {
		return stacked.StackTrace()
	}
	return errors.StackTrace{}
}
