package framework

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/financial_advisor/app/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	root(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"version": "v1",
		"name": "Financial Advisor"
	}`, w.Body.String())
}

func TestHealthz(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		isShuttingDown bool
		expectedStatus int
		expectedJSON   string
	}{
		{
			name:           "healthy status",
			isShuttingDown: false,
			expectedStatus: http.StatusOK,
			expectedJSON:   `{"status": "OK"}`,
		},
		{
			name:           "shutting down status",
			isShuttingDown: true,
			expectedStatus: http.StatusServiceUnavailable,
			expectedJSON:   `{"status": "SHUTTING_DOWN"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Set the shutdown state
			config.IsShuttingDown.Store(tt.isShuttingDown)
			defer config.IsShuttingDown.Store(false) // Reset after test

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/api/healthz", nil)

			healthz(ctx)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedJSON, w.Body.String())
		})
	}
}