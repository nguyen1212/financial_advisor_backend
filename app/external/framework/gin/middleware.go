package framework

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

func Secure() gin.HandlerFunc {
	return secure.New(secure.Config{
		FrameDeny:             true,
		ContentSecurityPolicy: "frame-ancestors 'none'",
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	})
}

// Headers add secure headers
func Headers(c *gin.Context) {
	c.Writer.Header().Add("Cache-Control", "no-store")

	c.Next()
}

// CorsMiddleware return the middleware instance
func CORS(allowOriginHosts []string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: allowOriginHosts,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Accept",
			"Authorization",
			"Cache-Control",
			"Content-Length",
			"Content-Type",
			"Cookie",
			"Origin",
			"Pragma",
			"X-Csrf-Token",
			"X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
