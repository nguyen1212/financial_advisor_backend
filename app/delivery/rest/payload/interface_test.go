package payload

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldBindWith(t *testing.T) {
	t.Parallel()

	type TestPayload struct {
		Name  string `json:"name" form:"name"`
		Email string `json:"email" form:"email"`
	}

	tests := []struct {
		name        string
		setupReq    func() *http.Request
		binding     binding.Binding
		wantPayload TestPayload
		wantErr     bool
	}{
		{
			name: "bind JSON payload successfully",
			setupReq: func() *http.Request {
				payload := TestPayload{Name: "John Doe", Email: "john@example.com"}
				jsonData, _ := json.Marshal(payload)
				req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			binding:     binding.JSON,
			wantPayload: TestPayload{Name: "John Doe", Email: "john@example.com"},
			wantErr:     false,
		},
		{
			name: "bind form data successfully",
			setupReq: func() *http.Request {
				form := "name=Jane+Doe&email=jane%40example.com"
				req, _ := http.NewRequest("POST", "/test", strings.NewReader(form))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				return req
			},
			binding:     binding.Form,
			wantPayload: TestPayload{Name: "Jane Doe", Email: "jane@example.com"},
			wantErr:     false,
		},
		{
			name: "bind query parameters successfully",
			setupReq: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test?name=Bob&email=bob@example.com", nil)
				return req
			},
			binding:     binding.Query,
			wantPayload: TestPayload{Name: "Bob", Email: "bob@example.com"},
			wantErr:     false,
		},
		{
			name: "bind invalid JSON should return error",
			setupReq: func() *http.Request {
				req, _ := http.NewRequest("POST", "/test", strings.NewReader("invalid json"))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			binding:     binding.JSON,
			wantPayload: TestPayload{},
			wantErr:     true,
		},
		{
			name: "bind empty request body with JSON binding",
			setupReq: func() *http.Request {
				req, _ := http.NewRequest("POST", "/test", nil)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			binding:     binding.JSON,
			wantPayload: TestPayload{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := tt.setupReq()
			var payload TestPayload

			err := ShouldBindWith(req, &payload, tt.binding)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantPayload, payload)
			}
		})
	}
}