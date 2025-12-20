package framework

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Recovery
// a custom recovery for Gin requests
// We decide to not use gin.Recovery() because it may log request headers into stderr under somce circumstances,
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Error(err.(error)) //nolint:errcheck,forcetypeassert

				c.AbortWithStatus(http.StatusInternalServerError)

				logrus.
					WithField("url", c.Request.URL.String()).
					WithField("method", c.Request.Method).
					WithField("user_agent", c.Request.UserAgent()).
					WithField("ip", c.ClientIP()).
					Errorln("recovered from panic")
			}
		}()

		c.Next()
	}
}
