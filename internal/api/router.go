package api

import (
	"b3-challenge/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures and returns the API router
func SetupRouter(tradeService service.TradeServiceInterface) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(LoggerMiddleware())
	router.Use(ErrorHandlerMiddleware())
	router.Use(CORSMiddleware())
	router.Use(RateLimitMiddleware())

	// Create handlers
	handlers := NewHandlers(tradeService)

	// Health check endpoint
	router.GET("/health", handlers.HealthCheck)

	// Swagger docs endpoint
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Tickers endpoints
		tickers := v1.Group("/tickers")
		{
			tickers.GET("", handlers.GetTickers)                               // GET /api/v1/tickers
			tickers.GET("/:ticker/aggregation", handlers.GetTickerAggregation) // GET /api/v1/tickers/{ticker}/aggregation
		}

		// Trades endpoints
		trades := v1.Group("/trades")
		{
			trades.GET("", handlers.GetTradesByDateRange) // GET /api/v1/trades?from=2024-01-01&to=2024-01-31
		}
	}

	// API documentation endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "B3 Challenge API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"health":                      "GET /health",
				"tickers":                     "GET /api/v1/tickers",
				"ticker_aggregation":          "GET /api/v1/tickers/{ticker}/aggregation",
				"trades_by_date_range":        "GET /api/v1/trades?from=YYYY-MM-DD&to=YYYY-MM-DD",
				"trades_last_7_business_days": "GET /api/v1/trades",
			},
			"documentation": gin.H{
				"swagger": "/docs/index.html",
			},
		})
	})

	return router
}
