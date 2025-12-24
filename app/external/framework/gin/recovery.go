package framework

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Recovery is a custom recovery for Gin requests
// We decide to not use gin.Recovery() because it may log request headers into stderr under somce circumstances,
func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
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
