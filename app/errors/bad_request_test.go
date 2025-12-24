package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewErrorBadRequest(t *testing.T) {
	t.Parallel()
	assert.Equal(t, errorBadRequest{
		code:    ErrorCodeBadRequest,
		message: "Bad Request",
	}, NewErrorBadRequest(ErrorCodeBadRequest, "Bad Request"))
}

func Test_ErrorBadRequest_Error(t *testing.T) {
	t.Parallel()
	err := errorBadRequest{
		code:    ErrorCodeBadRequest,
		message: "Bad Request",
	}
	assert.Equal(t, fmt.Sprintf("Type: %v, \tCode: %v, \tMessage: %s",
		ErrorTypeBadRequest, ErrorCodeBadRequest, "Bad Request"), err.Error())
}

func Test_ErrorBadRequest_Code(t *testing.T) {
	t.Parallel()
	err := errorBadRequest{
		code:    ErrorCodeBadRequest,
		message: "Bad Request",
	}
	assert.Equal(t, ErrorCodeBadRequest, err.Code())
}

func Test_ErrorBadRequest_Message(t *testing.T) {
	t.Parallel()
	err := errorBadRequest{
		code:    ErrorCodeBadRequest,
		message: "Bad Request",
	}
	assert.Equal(t, "Bad Request", err.Message())
}

func Test_ErrorBadRequest_HTTPCode(t *testing.T) {
	t.Parallel()
	err := errorBadRequest{
		code:    ErrorCodeBadRequest,
		message: "Bad Request",
	}
	assert.Equal(t, 400, err.HTTPCode())
}
