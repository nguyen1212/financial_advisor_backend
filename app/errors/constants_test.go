package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SystemErrors_Error(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		errors   SystemErrors
		expected string
	}{
		{
			name: "single error",
			errors: SystemErrors{
				NewErrorBadRequest(ErrorCodeBadRequest, "Invalid input"),
			},
			expected: "Type: TYPE_BAD_REQUEST, \tCode: CODE_BAD_REQUEST, \tMessage: Invalid input",
		},
		{
			name: "multiple errors",
			errors: SystemErrors{
				NewErrorBadRequest(ErrorCodeBadRequest, "Invalid input"),
				NewErrorForbidden(ErrorCodeForbidden, "Access denied"),
			},
			expected: "Type: TYPE_BAD_REQUEST, \tCode: CODE_BAD_REQUEST, \tMessage: Invalid input\nType: TYPE_FORBIDDEN, \tCode: CODE_FORBIDDEN, \tMessage: Access denied",
		},
		{
			name:     "empty errors",
			errors:   SystemErrors{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.errors.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_ErrorConstants(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		errorVar error
		expected string
	}{
		{
			name:     "ErrInternalServer",
			errorVar: ErrInternalServer,
			expected: "internal server error",
		},
		{
			name:     "ErrNotFound",
			errorVar: ErrNotFound,
			expected: "resource not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, tt.errorVar.Error())
		})
	}
}

