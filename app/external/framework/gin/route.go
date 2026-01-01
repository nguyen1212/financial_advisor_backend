package framework

import (
	"net/http"

	"github.com/financial_advisor/app/config"
	_ "github.com/financial_advisor/app/delivery/rest/docs"
	"github.com/financial_advisor/app/delivery/rest/handler"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handler define mapping routes
// @title financial_advisor backend
// @version 1.0
// @description This is the project of financial advisor backend.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
func Handler() *gin.Engine {
	var (
		router           = gin.New()
		newsHandler      = newNewsHandlerWrapper(handler.NewNewsHandler())
		publisherHandler = newPublisherHandlerWrapper(handler.NewPublisherHandler())
	)

	router.Use(
		Recovery,
		Secure(),
		Headers,
		CORS(config.Get().Cors),
		gin.Logger(),
	)

	if config.Get().ENV == "development" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	router.GET("/", root)
	router.GET("/api/healthz", healthz)

	// News routes
	router.GET("/api/v1/news", newsHandler.Find)
	router.GET("/api/v1/news/:id", newsHandler.Get)
	router.POST("/api/v1/news", newsHandler.Create)
	router.DELETE("/api/v1/news/:id", newsHandler.Delete)
	router.GET("/api/v1/news/search/suggestions", newsHandler.GetSearchSuggestions)
	router.GET("/api/v1/news/search", newsHandler.Search)

	// Publisher routes
	router.GET("/api/v1/publishers", publisherHandler.Find)
	router.GET("/api/v1/publishers/:id", publisherHandler.Get)
	router.POST("/api/v1/publishers", publisherHandler.Create)

	return router
}

func root(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, struct {
		Version string `json:"version"`
		Name    string `json:"name"`
	}{
		Version: "v1",
		Name:    "Financial Advisor",
	})
}

func healthz(c *gin.Context) {
	if config.IsShuttingDown.Load() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "SHUTTING_DOWN",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}
