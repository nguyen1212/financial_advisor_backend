package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewErrorConflicted(t *testing.T) {
	t.Parallel()

	assert.Equal(t, errorConflicted{
		code:    ErrorCodeConflicted,
		message: "Resource Conflict",
	}, NewErrorConflicted(ErrorCodeConflicted, "Resource Conflict"))
}

func Test_ErrorConflicted_Error(t *testing.T) {
	t.Parallel()

	err := errorConflicted{
		code:    ErrorCodeConflicted,
		message: "Resource Conflict",
	}
	assert.Equal(t, fmt.Sprintf("Type: %v, \tCode: %v, \tMessage: %s",
		ErrorTypeConflicted, ErrorCodeConflicted, "Resource Conflict"), err.Error())
}

func Test_ErrorConflicted_Code(t *testing.T) {
	t.Parallel()

	err := errorConflicted{
		code:    ErrorCodeConflicted,
		message: "Resource Conflict",
	}
	assert.Equal(t, ErrorCodeConflicted, err.Code())
}

func Test_ErrorConflicted_Message(t *testing.T) {
	t.Parallel()

	err := errorConflicted{
		code:    ErrorCodeConflicted,
		message: "Resource Conflict",
	}
	assert.Equal(t, "Resource Conflict", err.Message())
}

func Test_ErrorConflicted_HTTPCode(t *testing.T) {
	t.Parallel()

	err := errorConflicted{
		code:    ErrorCodeConflicted,
		message: "Resource Conflict",
	}
	assert.Equal(t, 409, err.HTTPCode())
}

func Test_ErrorConflicted_Type(t *testing.T) {
	t.Parallel()

	err := errorConflicted{
		code:    ErrorCodeConflicted,
		message: "Resource Conflict",
	}
	assert.Equal(t, ErrorTypeConflicted, err.Type())
}