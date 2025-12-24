package framework

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/financial_advisor/app/delivery/rest/handler/mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_baseHandlerWrapper_Find(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	t.Run("find is not implemented", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		hdl := &baseHandlerWrapper{}

		hdl.Find(ginCtx)

		assert.Equal(t, http.StatusNotImplemented, ginCtx.Writer.Status())
		assert.JSONEq(t, `{
      "errors": [
        {
        "code":"CODE_NOT_IMPLEMENTED",
        "message": "method find is not implemented",
        "type":"TYPE_NOT_IMPLEMENTED"
        }
      ]
		}`, w.Body.String())
	})

	t.Run("inner find handler failed", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
		mockHandler := mock.NewMockFindHandler(mockCtrl)
		hdl := &baseHandlerWrapper{
			findHandler: mockHandler,
		}

		mockHandler.EXPECT().Find(ginCtx.Request).Return(nil, nil, errors.New("find error"))

		hdl.Find(ginCtx)

		assert.Equal(t, http.StatusInternalServerError, ginCtx.Writer.Status())
		assert.JSONEq(t, `{
      "errors": [
        {
          "code":"CODE_INTERNAL",
          "message": "internal server error",
          "type":"TYPE_INTERNAL"
        }
      ]
		}`, w.Body.String())
	})
}

func Test_baseHandlerWrapper_Create(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	t.Run("create is not implemented", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
		hdl := &baseHandlerWrapper{}

		hdl.Create(ginCtx)

		assert.Equal(t, http.StatusNotImplemented, ginCtx.Writer.Status())
		assert.JSONEq(t, `{
      "errors": [
        {
          "code":"CODE_NOT_IMPLEMENTED",
          "message": "method create is not implemented",
          "type":"TYPE_NOT_IMPLEMENTED"
        }
      ]
		}`, w.Body.String())
	})

	t.Run("inner create handler failed", func(t *testing.T) {
		t.Parallel()

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		w := httptest.NewRecorder()
		ginCtx, _ := gin.CreateTestContext(w)
		ginCtx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
		mockHandler := mock.NewMockCreateHandler(mockCtrl)
		hdl := &baseHandlerWrapper{
			createHandler: mockHandler,
		}

		mockHandler.EXPECT().Create(ginCtx.Request).Return(nil, errors.New("create error"))

		hdl.Create(ginCtx)

		assert.Equal(t, http.StatusInternalServerError, ginCtx.Writer.Status())
		assert.JSONEq(t, `{
      "errors": [
        {
          "code":"CODE_INTERNAL",
          "message": "internal server error",
          "type":"TYPE_INTERNAL"
        }
      ]
		}`, w.Body.String())
	})
}
