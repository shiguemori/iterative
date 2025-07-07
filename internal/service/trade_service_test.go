package service

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"

	"b3-challenge/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisCache struct {
	mock.Mock
}

func (m *MockRedisCache) GetString(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisCache) SetJSON(key string, value interface{}, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

// MockTradeRepository is a mock implementation of TradeRepository
type MockTradeRepository struct {
	mock.Mock
}

func (m *MockTradeRepository) CreateTrades(ctx context.Context, trades []models.Trade) error {
	args := m.Called(ctx, trades)
	return args.Error(0)
}

func (m *MockTradeRepository) GetTradeAggregation(ctx context.Context, ticker string, fromDate *time.Time) (*models.TradeAggregation, error) {
	args := m.Called(ctx, ticker, fromDate)
	return args.Get(0).(*models.TradeAggregation), args.Error(1)
}

func (m *MockTradeRepository) GetTradesByDateRange(ctx context.Context, fromDate, toDate time.Time) ([]models.Trade, error) {
	args := m.Called(ctx, fromDate, toDate)
	return args.Get(0).([]models.Trade), args.Error(1)
}

func (m *MockTradeRepository) DeleteOldTrades(ctx context.Context, beforeDate time.Time) error {
	args := m.Called(ctx, beforeDate)
	return args.Error(0)
}

func (m *MockTradeRepository) GetDistinctTickers(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	return args.Get(0).([]string), args.Error(1)
}

func TestTradeService_GetTradeAggregation(t *testing.T) {
	// Setup
	mockRedis := new(MockRedisCache)
	mockRepo := new(MockTradeRepository)

	service := &TradeService{
		tradeRepo: mockRepo,
		cache:     mockRedis,
	}

	ctx := context.Background()
	ticker := "PETR4"

	expectedAggregation := &models.TradeAggregation{
		Ticker:         ticker,
		MaxRangeValue:  100.50,
		MaxDailyVolume: 1000000,
	}

	// cache miss
	cacheKey := fmt.Sprintf("GetTradeAggregation:%s", ticker)
	mockRedis.On("GetString", cacheKey).Return("", redis.Nil)

	// Mock expectations
	mockRepo.On("GetTradeAggregation", ctx, ticker, mock.AnythingOfType("*time.Time")).
		Return(expectedAggregation, nil)

	// set no cache
	mockRedis.On("SetJSON", cacheKey, mock.Anything, 10*time.Minute).Return(nil)

	// Execute
	result, err := service.GetTradeAggregation(ctx, ticker)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedAggregation, result)
	mockRepo.AssertExpectations(t)
}

func TestTradeService_ValidateTicker(t *testing.T) {
	// Setup
	mockRedis := new(MockRedisCache)
	mockRepo := new(MockTradeRepository)

	service := &TradeService{
		tradeRepo: mockRepo,
		cache:     mockRedis,
	}

	ctx := context.Background()

	tickers := []string{"PETR4", "VALE3", "ITUB4"}

	// Mock expectations
	mockRepo.On("GetDistinctTickers", ctx).Return(tickers, nil)

	// Test valid ticker
	exists, err := service.ValidateTicker(ctx, "PETR4")
	assert.NoError(t, err)
	assert.True(t, exists)

	// Test invalid ticker
	exists, err = service.ValidateTicker(ctx, "INVALID")
	assert.NoError(t, err)
	assert.False(t, exists)

	mockRepo.AssertExpectations(t)
}

func TestTradeService_GetBusinessDaysAgo(t *testing.T) {
	// Setup
	mockRedis := new(MockRedisCache)
	mockRepo := new(MockTradeRepository)

	service := &TradeService{
		tradeRepo: mockRepo,
		cache:     mockRedis,
	}

	// Test getting business days ago
	result := service.getBusinessDaysAgo(7)

	// Should be approximately 7 business days ago
	// This is a simple test - in practice you'd want more sophisticated date testing
	assert.True(t, result.Before(time.Now()))
}

func BenchmarkTradeService_GetTradeAggregation(b *testing.B) {
	// Setup
	mockRedis := new(MockRedisCache)
	mockRepo := new(MockTradeRepository)

	service := &TradeService{
		tradeRepo: mockRepo,
		cache:     mockRedis,
	}

	ctx := context.Background()
	ticker := "PETR4"

	expectedAggregation := &models.TradeAggregation{
		Ticker:         ticker,
		MaxRangeValue:  100.50,
		MaxDailyVolume: 1000000,
	}

	// cache miss
	cacheKey := fmt.Sprintf("GetTradeAggregation:%s", ticker)
	mockRedis.On("GetString", cacheKey).Return("", redis.Nil)

	// Mock expectations
	mockRepo.On("GetTradeAggregation", ctx, ticker, mock.AnythingOfType("*time.Time")).
		Return(expectedAggregation, nil)

	// set no cache
	mockRedis.On("SetJSON", cacheKey, mock.Anything, 10*time.Minute).Return(nil)

	// Benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetTradeAggregation(ctx, ticker)
		if err != nil {
			b.Fatal(err)
		}
	}
}
