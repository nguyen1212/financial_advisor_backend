package payload

import (
	"testing"

	appErrors "github.com/financial_advisor/app/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestValidateStruct struct {
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"min=18,max=100"`
}

type TestTransformStruct struct {
	Value string
}

type TestSliceStruct struct {
	Users []string `validate:"required,dive,min=2"`
}

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload interface{}
		wantErr bool
		errType string
		errMsg  string
	}{
		{
			name: "valid struct should pass validation",
			payload: TestValidateStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "missing required field should return validation error",
			payload: TestValidateStruct{
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: true,
			errType: "SystemErrors",
			errMsg:  "Invalid value for Name with tag required: ",
		},
		{
			name: "invalid email should return validation error",
			payload: TestValidateStruct{
				Name:  "John Doe",
				Email: "invalid-email",
				Age:   25,
			},
			wantErr: true,
			errType: "SystemErrors",
			errMsg:  "Invalid value for Email with tag email: invalid-email",
		},
		{
			name: "age below minimum should return validation error",
			payload: TestValidateStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   15,
			},
			wantErr: true,
			errType: "SystemErrors",
			errMsg:  "Invalid value for Age with tag min: 15",
		},
		{
			name: "age above maximum should return validation error",
			payload: TestValidateStruct{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   150,
			},
			wantErr: true,
			errType: "SystemErrors",
			errMsg:  "Invalid value for Age with tag max: 150",
		},
		{
			name: "multiple validation errors should return all errors",
			payload: TestValidateStruct{
				Email: "invalid-email",
				Age:   15,
			},
			wantErr: true,
			errType: "SystemErrors",
		},
		{
			name: "slice validation with invalid items",
			payload: TestSliceStruct{
				Users: []string{"a", "b"}, // both items are too short (min=2)
			},
			wantErr: true,
			errType: "SystemErrors",
			errMsg:  "Invalid value for Users with tag min: a",
		},
		{
			name:    "nil payload should return error",
			payload: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validate(tt.payload)

			if tt.wantErr {
				require.Error(t, err)

				if tt.errType == "SystemErrors" {
					var systemErrors appErrors.SystemErrors
					require.ErrorAs(t, err, &systemErrors)
					require.Greater(t, len(systemErrors), 0)

					if tt.errMsg != "" {
						found := false
						for _, sysErr := range systemErrors {
							if sysErr.Message() == tt.errMsg {
								found = true
								break
							}
						}
						assert.True(t, found, "Expected error message not found: %s", tt.errMsg)
					}
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNormalizeFieldName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "field without slice notation",
			input:    "Name",
			expected: "Name",
		},
		{
			name:     "field with slice notation",
			input:    "Users[0]",
			expected: "Users",
		},
		{
			name:     "field with nested slice notation",
			input:    "Users[0].Posts[1]",
			expected: "Users",
		},
		{
			name:     "field with complex slice notation",
			input:    "Categories[5]",
			expected: "Categories",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "field with brackets but not at start",
			input:    "Name[test]other",
			expected: "Name",
		},
		{
			name:     "field starting with bracket",
			input:    "[0]Name",
			expected: "[0]Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := normalizeFieldName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTransform(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		payload interface{}
		wantErr bool
	}{
		{
			name: "transform valid struct",
			payload: &TestTransformStruct{
				Value: "test value",
			},
			wantErr: false,
		},
		{
			name:    "transform nil payload should return error",
			payload: nil,
			wantErr: true,
		},
		{
			name: "transform empty struct",
			payload: &TestTransformStruct{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := transform(tt.payload)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestValidateIntegration tests the validate function behavior with actual validation tags
func TestValidateIntegration(t *testing.T) {
	t.Parallel()

	type IntegrationTestStruct struct {
		RequiredField string `validate:"required"`
		EmailField    string `validate:"email"`
		URLField      string `validate:"url"`
		OneOfField    string `validate:"oneof=finance military sports"`
	}

	tests := []struct {
		name        string
		payload     IntegrationTestStruct
		expectError bool
		errorCount  int
	}{
		{
			name: "all valid fields",
			payload: IntegrationTestStruct{
				RequiredField: "test",
				EmailField:    "test@example.com",
				URLField:      "https://example.com",
				OneOfField:    "finance",
			},
			expectError: false,
		},
		{
			name: "invalid email and oneof",
			payload: IntegrationTestStruct{
				RequiredField: "test",
				EmailField:    "invalid-email",
				URLField:      "https://example.com",
				OneOfField:    "invalid-category",
			},
			expectError: true,
			errorCount:  2,
		},
		{
			name: "missing required field",
			payload: IntegrationTestStruct{
				EmailField: "test@example.com",
				URLField:   "https://example.com",
				OneOfField: "finance",
			},
			expectError: true,
			errorCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validate(tt.payload)

			if tt.expectError {
				require.Error(t, err)
				var systemErrors appErrors.SystemErrors
				require.ErrorAs(t, err, &systemErrors)
				assert.Equal(t, tt.errorCount, len(systemErrors))
			} else {
				require.NoError(t, err)
			}
		})
	}
}