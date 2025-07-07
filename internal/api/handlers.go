package api

import (
	"net/http"
	"strings"
	"time"

	"b3-challenge/internal/service"

	"github.com/gin-gonic/gin"
)

// Handlers contains all API handlers
type Handlers struct {
	tradeService service.TradeServiceInterface
}

// NewHandlers creates a new handlers instance
func NewHandlers(tradeService service.TradeServiceInterface) *Handlers {
	return &Handlers{
		tradeService: tradeService,
	}
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// HealthCheck godoc
// @Summary Health check endpoint
// @Tags Health
// @Produce json
// @Success 200 {object} APIResponse
// @Router /health [get]
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: "B3 Challenge API is running",
		Data: gin.H{
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		},
	})
}

// GetTickerAggregation godoc
// @Summary Get aggregated data for a specific ticker
// @Tags Tickers
// @Produce json
// @Param ticker path string true "Ticker symbol (e.g. PETR4, VALE3)"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /tickers/{ticker}/aggregation [get]
func (h *Handlers) GetTickerAggregation(c *gin.Context) {
	ticker := strings.ToUpper(c.Param("ticker"))

	if ticker == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "ticker parameter is required",
		})
		return
	}

	// Validate ticker exists
	exists, err := h.tradeService.ValidateTicker(c.Request.Context(), ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "failed to validate ticker",
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Error:   "ticker not found",
		})
		return
	}

	aggregation, err := h.tradeService.GetTradeAggregation(c.Request.Context(), ticker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "failed to get ticker aggregation",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    aggregation,
	})
}

// GetTradesByDateRange godoc
// @Summary Get trades within a specific date range
// @Tags Trades
// @Produce json
// @Param from query string false "Start date in YYYY-MM-DD format"
// @Param to query string false "End date in YYYY-MM-DD format"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /trades [get]
func (h *Handlers) GetTradesByDateRange(c *gin.Context) {

	fromStr := c.Query("from")
	toStr := c.Query("to")

	var fromDate, toDate time.Time
	var err error

	if fromStr == "" || toStr == "" {
		// Default to last 7 business days
		trades, err := h.tradeService.GetLast7BusinessDaysTrades(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, APIResponse{
				Success: false,
				Error:   "failed to get trades",
			})
			return
		}

		c.JSON(http.StatusOK, APIResponse{
			Success: true,
			Data:    trades,
			Message: "Trades from last 7 business days",
		})
		return
	}

	// Parse dates
	fromDate, err = time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "invalid 'from' date format. Use YYYY-MM-DD",
		})
		return
	}

	toDate, err = time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "invalid 'to' date format. Use YYYY-MM-DD",
		})
		return
	}

	// Validate date range
	if fromDate.After(toDate) {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "'from' date cannot be after 'to' date",
		})
		return
	}

	// Limit range to prevent excessive queries
	if toDate.Sub(fromDate).Hours() > 24*90 { // 90 days limit
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "date range cannot exceed 90 days",
		})
		return
	}

	trades, err := h.tradeService.GetTradesByDateRange(c.Request.Context(), fromDate, toDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "failed to get trades by date range",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    trades,
		Message: "Trades from specified date range",
	})
}

// GetTickers godoc
// @Summary list all available tickers
// @Tags Tickers
// @Produce json
// @Success 200 {array} string
// @Failure 500 {object} error
// @Router /tickers [get]
func (h *Handlers) GetTickers(c *gin.Context) {
	tickers, err := h.tradeService.GetDistinctTickers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   "failed to get tickers",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    tickers,
		Message: "Available tickers",
	})
}
