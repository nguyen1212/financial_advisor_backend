package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewErrorNotImplemented(t *testing.T) {
	t.Parallel()
	assert.Equal(t, errorNotImplemented{
		code:    ErrorCodeNotImplemented,
		message: "Feature Not Implemented",
	}, NewErrorNotImplemented(ErrorCodeNotImplemented, "Feature Not Implemented"))
}

func Test_ErrorNotImplemented_Error(t *testing.T) {
	t.Parallel()
	err := errorNotImplemented{
		code:    ErrorCodeNotImplemented,
		message: "Feature Not Implemented",
	}
	assert.Equal(t, fmt.Sprintf("Type: %v, \tCode: %v, \tMessage: %s",
		ErrorTypeNotImplemented, ErrorCodeNotImplemented, "Feature Not Implemented"), err.Error())
}

func Test_ErrorNotImplemented_Code(t *testing.T) {
	t.Parallel()
	err := errorNotImplemented{
		code:    ErrorCodeNotImplemented,
		message: "Feature Not Implemented",
	}
	assert.Equal(t, ErrorCodeNotImplemented, err.Code())
}

func Test_ErrorNotImplemented_Message(t *testing.T) {
	t.Parallel()
	err := errorNotImplemented{
		code:    ErrorCodeNotImplemented,
		message: "Feature Not Implemented",
	}
	assert.Equal(t, "Feature Not Implemented", err.Message())
}

func Test_ErrorNotImplemented_HTTPCode(t *testing.T) {
	t.Parallel()
	err := errorNotImplemented{
		code:    ErrorCodeNotImplemented,
		message: "Feature Not Implemented",
	}
	assert.Equal(t, 501, err.HTTPCode())
}

func Test_ErrorNotImplemented_Type(t *testing.T) {
	t.Parallel()
	err := errorNotImplemented{
		code:    ErrorCodeNotImplemented,
		message: "Feature Not Implemented",
	}
	assert.Equal(t, ErrorTypeNotImplemented, err.Type())
}