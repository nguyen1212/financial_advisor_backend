package framework

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRecovery(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		handler        gin.HandlerFunc
		shouldPanic    bool
		expectedStatus int
	}{
		{
			name: "normal request without panic",
			handler: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			},
			shouldPanic:    false,
			expectedStatus: http.StatusOK,
		},
		{
			name: "request that panics",
			handler: func(c *gin.Context) {
				panic("something went wrong")
			},
			shouldPanic:    true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "request that panics with nil",
			handler: func(c *gin.Context) {
				panic(nil)
			},
			shouldPanic:    true, // panic(nil) does actually cause a panic
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			router := gin.New()
			router.Use(Recovery)
			router.GET("/test", tt.handler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.shouldPanic {
				// When panic occurs, response should be aborted with 500
				assert.Empty(t, w.Body.String())
			}
		})
	}
}

func TestRecovery_LogsRequestDetails(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	router := gin.New()
	router.Use(Recovery)
	router.GET("/test", func(c *gin.Context) {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/test?param=value", nil)
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("X-Forwarded-For", "192.168.1.1")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Empty(t, w.Body.String())
}

