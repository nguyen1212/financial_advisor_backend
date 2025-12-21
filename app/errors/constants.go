// Package errors define error of the application
package errors

import (
	"errors"
	"strings"
)

var ErrInternalServer = errors.New("internal server error")

type ErrorType string

const (
	ErrorTypeInternal   ErrorType = "TYPE_INTERNAL"
	ErrorTypeBadRequest ErrorType = "TYPE_BAD_REQUEST"
)

type ErrorCode string

// general error codes
const (
	ErrorCodeInternal   ErrorCode = "CODE_INTERNAL"
	ErrorCodeBadRequest ErrorCode = "CODE_BAD_REQUEST"
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
