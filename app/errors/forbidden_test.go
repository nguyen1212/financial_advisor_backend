package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewErrorForbidden(t *testing.T) {
	t.Parallel()
	assert.Equal(t, errorForbidden{
		code:    ErrorCodeForbidden,
		message: "Access Forbidden",
	}, NewErrorForbidden(ErrorCodeForbidden, "Access Forbidden"))
}

func Test_ErrorForbidden_Error(t *testing.T) {
	t.Parallel()
	err := errorForbidden{
		code:    ErrorCodeForbidden,
		message: "Access Forbidden",
	}
	assert.Equal(t, fmt.Sprintf("Type: %v, \tCode: %v, \tMessage: %s",
		ErrorTypeForbidden, ErrorCodeForbidden, "Access Forbidden"), err.Error())
}

func Test_ErrorForbidden_Code(t *testing.T) {
	t.Parallel()
	err := errorForbidden{
		code:    ErrorCodeForbidden,
		message: "Access Forbidden",
	}
	assert.Equal(t, ErrorCodeForbidden, err.Code())
}

func Test_ErrorForbidden_Message(t *testing.T) {
	t.Parallel()
	err := errorForbidden{
		code:    ErrorCodeForbidden,
		message: "Access Forbidden",
	}
	assert.Equal(t, "Access Forbidden", err.Message())
}

func Test_ErrorForbidden_HTTPCode(t *testing.T) {
	t.Parallel()
	err := errorForbidden{
		code:    ErrorCodeForbidden,
		message: "Access Forbidden",
	}
	assert.Equal(t, 403, err.HTTPCode())
}

func Test_ErrorForbidden_Type(t *testing.T) {
	t.Parallel()
	err := errorForbidden{
		code:    ErrorCodeForbidden,
		message: "Access Forbidden",
	}
	assert.Equal(t, ErrorTypeForbidden, err.Type())
}