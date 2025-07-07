package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"b3-challenge/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTradeService is a mock implementation of TradeService
type MockTradeService struct {
	mock.Mock
}

func (m *MockTradeService) GetTradeAggregation(ctx context.Context, ticker string) (*models.TradeAggregation, error) {
	args := m.Called(ctx, ticker)
	return args.Get(0).(*models.TradeAggregation), args.Error(1)
}

func (m *MockTradeService) GetAllTickersAggregation(ctx context.Context) ([]models.TradeAggregation, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.TradeAggregation), args.Error(1)
}

func (m *MockTradeService) GetTradesByDateRange(ctx context.Context, fromDate, toDate time.Time) ([]models.Trade, error) {
	args := m.Called(ctx, fromDate, toDate)
	return args.Get(0).([]models.Trade), args.Error(1)
}

func (m *MockTradeService) GetLast7BusinessDaysTrades(ctx context.Context) ([]models.Trade, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Trade), args.Error(1)
}

func (m *MockTradeService) ValidateTicker(ctx context.Context, ticker string) (bool, error) {
	args := m.Called(ctx, ticker)
	return args.Bool(0), args.Error(1)
}

func (m *MockTradeService) GetDistinctTickers(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func setupTestRouter(mockService *MockTradeService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handlers := NewHandlers(mockService)

	router.GET("/health", handlers.HealthCheck)
	router.GET("/api/v1/tickers", handlers.GetTickers)
	router.GET("/api/v1/tickers/:ticker/aggregation", handlers.GetTickerAggregation)
	router.GET("/api/v1/trades", handlers.GetTradesByDateRange)

	return router
}

func TestHealthCheck(t *testing.T) {
	mockService := new(MockTradeService)
	router := setupTestRouter(mockService)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "B3 Challenge API is running", response.Message)
}

func TestGetTickers(t *testing.T) {
	mockService := new(MockTradeService)
	router := setupTestRouter(mockService)

	expectedTickers := []string{"PETR4", "VALE3", "ITUB4"}
	mockService.On("GetDistinctTickers", mock.Anything).Return(expectedTickers, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tickers", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Available tickers", response.Message)

	// Convert data to slice of strings
	dataBytes, _ := json.Marshal(response.Data)
	var tickers []string
	err = json.Unmarshal(dataBytes, &tickers)
	if err != nil {
		t.Errorf("Failed to unmarshal tickers: %v", err)
		return
	}
	assert.Equal(t, expectedTickers, tickers)

	mockService.AssertExpectations(t)
}

func TestGetTickerAggregation_Success(t *testing.T) {
	mockService := new(MockTradeService)
	router := setupTestRouter(mockService)

	ticker := "PETR4"
	expectedAggregation := &models.TradeAggregation{
		Ticker:         ticker,
		MaxRangeValue:  100.50,
		MaxDailyVolume: 1000000,
	}

	mockService.On("ValidateTicker", mock.Anything, ticker).Return(true, nil)
	mockService.On("GetTradeAggregation", mock.Anything, ticker).Return(expectedAggregation, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tickers/PETR4/aggregation", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)

	mockService.AssertExpectations(t)
}

func TestGetTickerAggregation_NotFound(t *testing.T) {
	mockService := new(MockTradeService)
	router := setupTestRouter(mockService)

	ticker := "INVALID"
	mockService.On("ValidateTicker", mock.Anything, ticker).Return(false, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/tickers/INVALID/aggregation", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "ticker not found", response.Error)

	mockService.AssertExpectations(t)
}

func TestGetTradesByDateRange_WithDates(t *testing.T) {
	mockService := new(MockTradeService)
	router := setupTestRouter(mockService)

	fromDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	expectedTrades := []models.Trade{
		{ID: 1, CodigoInstrumento: "PETR4", PrecoNegocio: 100.50, QuantidadeNegociada: 1000},
	}

	mockService.On("GetTradesByDateRange", mock.Anything, fromDate, toDate).Return(expectedTrades, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/trades?from=2024-01-01&to=2024-01-31", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Trades from specified date range", response.Message)

	mockService.AssertExpectations(t)
}

func TestGetTradesByDateRange_Last7Days(t *testing.T) {
	mockService := new(MockTradeService)
	router := setupTestRouter(mockService)

	expectedTrades := []models.Trade{
		{ID: 1, CodigoInstrumento: "PETR4", PrecoNegocio: 100.50, QuantidadeNegociada: 1000},
	}

	mockService.On("GetLast7BusinessDaysTrades", mock.Anything).Return(expectedTrades, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/trades", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.Equal(t, "Trades from last 7 business days", response.Message)

	mockService.AssertExpectations(t)
}

func BenchmarkGetTickerAggregation(b *testing.B) {
	mockService := new(MockTradeService)
	router := setupTestRouter(mockService)

	ticker := "PETR4"
	expectedAggregation := &models.TradeAggregation{
		Ticker:         ticker,
		MaxRangeValue:  100.50,
		MaxDailyVolume: 1000000,
	}

	mockService.On("ValidateTicker", mock.Anything, ticker).Return(true, nil)
	mockService.On("GetTradeAggregation", mock.Anything, ticker).Return(expectedAggregation, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/tickers/PETR4/aggregation", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatal("Expected status 200")
		}
	}
}
