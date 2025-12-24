package errors

import (
	"fmt"
	"net/http"
)

type errorConflicted struct {
	code    ErrorCode
	message string
}

func NewErrorConflicted(code ErrorCode, message string) SystemError {
	return errorConflicted{
		code:    code,
		message: message,
	}
}

func (e errorConflicted) Type() ErrorType {
	return ErrorTypeConflicted
}

func (e errorConflicted) Code() ErrorCode {
	return e.code
}

func (e errorConflicted) Message() string {
	return e.message
}

func (e errorConflicted) HTTPCode() int {
	return http.StatusConflict
}

func (e errorConflicted) Error() string {
	return fmt.Sprintf("Type: %v, \tCode: %v, \tMessage: %s",
		ErrorTypeConflicted, e.Code(), e.Message())
}
