// Package errors define error of the application
package errors

import (
	"errors"
	"strings"
)

var (
	ErrInternalServer = errors.New("internal server error")
	ErrNotFound       = errors.New("resource not found")
	ErrConflicted     = errors.New("resource conflict")
)

type ErrorType string

const (
	ErrorTypeInternal       ErrorType = "TYPE_INTERNAL"
	ErrorTypeBadRequest     ErrorType = "TYPE_BAD_REQUEST"
	ErrorTypeForbidden      ErrorType = "TYPE_FORBIDDEN"
	ErrorTypeNotImplemented ErrorType = "TYPE_NOT_IMPLEMENTED"
	ErrorTypeConflicted     ErrorType = "TYPE_CONFLICTED"
)

type ErrorCode string

// general error codes
const (
	ErrorCodeInternal       ErrorCode = "CODE_INTERNAL"
	ErrorCodeBadRequest     ErrorCode = "CODE_BAD_REQUEST"
	ErrorCodeForbidden      ErrorCode = "CODE_FORBIDDEN"
	ErrorCodeNotImplemented ErrorCode = "CODE_NOT_IMPLEMENTED"
	ErrorCodeConflicted     ErrorCode = "CODE_CONFLICTED"
)

// application error codes
const (
	ErrorCodePublisherNotFound ErrorCode = "CODE_PUBLISHER_NOT_FOUND"
	ErrorCodeNewsNotFound      ErrorCode = "CODE_NEWS_NOT_FOUND"
	ErrorCodeURLTooLong        ErrorCode = "CODE_URL_TOO_LONG"
	ErrorCodeURLInvalid        ErrorCode = "CODE_URL_INVALID"
)

type SystemError interface {
	Type() ErrorType
	Code() ErrorCode
	Message() string
	HTTPCode() int
	Error() string
}

// SystemErrors an array of system errors
type SystemErrors []SystemError

// Error implements error interface
func (errs SystemErrors) Error() string {
	errsString := make([]string, 0, len(errs))
	for _, err := range errs {
		errsString = append(errsString, err.Error())
	}

	return strings.Join(errsString, "\n")
}
