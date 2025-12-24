package framework

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	appError "github.com/financial_advisor/app/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRenderData(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		data         any
		paging       any
		expectedJSON string
	}{
		{
			name:         "render data with paging",
			data:         map[string]string{"test": "value"},
			paging:       map[string]int{"page": 1, "size": 10},
			expectedJSON: `{"data":{"test":"value"},"paging":{"page":1,"size":10}}`,
		},
		{
			name:         "render data without paging",
			data:         []string{"item1", "item2"},
			paging:       nil,
			expectedJSON: `{"data":["item1","item2"]}`,
		},
		{
			name:         "render nil data",
			data:         nil,
			paging:       nil,
			expectedJSON: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			RenderData(ctx, tt.data, tt.paging)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, tt.expectedJSON, w.Body.String())
		})
	}
}

func TestRenderErrors(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedJSON   string
	}{
		{
			name: "system errors (bad request)",
			err: appError.SystemErrors{
				appError.NewErrorBadRequest(appError.ErrorCodeBadRequest, "Invalid input"),
				appError.NewErrorBadRequest(appError.ErrorCodeBadRequest, "Missing field"),
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: `{
				"errors": [
					{
						"type": "TYPE_BAD_REQUEST",
						"code": "CODE_BAD_REQUEST",
						"message": "Invalid input"
					},
					{
						"type": "TYPE_BAD_REQUEST",
						"code": "CODE_BAD_REQUEST",
						"message": "Missing field"
					}
				]
			}`,
		},
		{
			name:           "single system error (forbidden)",
			err:            appError.NewErrorForbidden(appError.ErrorCodeForbidden, "Access denied"),
			expectedStatus: http.StatusForbidden,
			expectedJSON: `{
				"errors": [
					{
						"type": "TYPE_FORBIDDEN",
						"code": "CODE_FORBIDDEN",
						"message": "Access denied"
					}
				]
			}`,
		},
		{
			name:           "single system error (not implemented)",
			err:            appError.NewErrorNotImplemented(appError.ErrorCodeNotImplemented, "Feature not available"),
			expectedStatus: http.StatusNotImplemented,
			expectedJSON: `{
				"errors": [
					{
						"type": "TYPE_NOT_IMPLEMENTED",
						"code": "CODE_NOT_IMPLEMENTED",
						"message": "Feature not available"
					}
				]
			}`,
		},
		{
			name:           "generic error (internal server error)",
			err:            errors.New("unexpected error"),
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: `{
				"errors": [
					{
						"type": "TYPE_INTERNAL",
						"code": "CODE_INTERNAL",
						"message": "internal server error"
					}
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			RenderErrors(ctx, tt.err)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedJSON, w.Body.String())
		})
	}
}

func Test_fromSystemErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		errs     appError.SystemErrors
		expected []responseError
	}{
		{
			name: "convert system errors",
			errs: appError.SystemErrors{
				appError.NewErrorBadRequest(appError.ErrorCodeBadRequest, "Bad input"),
				appError.NewErrorForbidden(appError.ErrorCodeForbidden, "Access denied"),
			},
			expected: []responseError{
				{
					Type:    "TYPE_BAD_REQUEST",
					Code:    "CODE_BAD_REQUEST",
					Message: "Bad input",
				},
				{
					Type:    "TYPE_FORBIDDEN",
					Code:    "CODE_FORBIDDEN",
					Message: "Access denied",
				},
			},
		},
		{
			name:     "empty system errors",
			errs:     appError.SystemErrors{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := fromSystemErrors(tt.errs)
			assert.Equal(t, tt.expected, result)
		})
	}
}