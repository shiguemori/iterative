package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"b3-challenge/internal/models"
	"b3-challenge/internal/repository"

	redis "b3-challenge/pkg/cache"
)

// TradeService handles business logic for trades
type TradeService struct {
	tradeRepo repository.TradeRepositoryInterface
	cache     redis.Cache
}

// NewTradeService creates a new trade service
func NewTradeService(tradeRepo repository.TradeRepositoryInterface, cache *redis.RedisCache) *TradeService {
	return &TradeService{
		tradeRepo: tradeRepo,
		cache:     cache,
	}
}

// GetTradeAggregation returns aggregated trade data for a ticker
func (s *TradeService) GetTradeAggregation(ctx context.Context, ticker string) (*models.TradeAggregation, error) {
	cacheKey := fmt.Sprintf("GetTradeAggregation:%s", ticker)

	cached, err := s.cache.GetString(cacheKey)
	if err == nil && cached != "" {
		var result models.TradeAggregation
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &result, nil
		}
	}

	// Calculate date 7 business days ago
	fromDate := s.getBusinessDaysAgo(7)

	aggregation, err := s.tradeRepo.GetTradeAggregation(ctx, ticker, &fromDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get trade aggregation for ticker %s: %w", ticker, err)
	}

	jsonData, _ := json.Marshal(aggregation)
	_ = s.cache.SetJSON(cacheKey, jsonData, 10*time.Minute)

	return aggregation, nil
}

// GetTradesByDateRange returns trades within a specific date range
func (s *TradeService) GetTradesByDateRange(ctx context.Context, fromDate, toDate time.Time) ([]models.Trade, error) {
	cacheKey := fmt.Sprintf("GetTradesByDateRange:%s:%s", fromDate, toDate)

	cached, err := s.cache.GetString(cacheKey)
	if err == nil && cached != "" {
		var result []models.Trade
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result, nil
		}
	}

	trades, err := s.tradeRepo.GetTradesByDateRange(ctx, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get trades by date range: %w", err)
	}

	jsonData, _ := json.Marshal(trades)
	_ = s.cache.SetJSON(cacheKey, jsonData, 10*time.Minute)

	return trades, nil
}

// GetLast7BusinessDaysTrades returns trades from the last 7 business days
func (s *TradeService) GetLast7BusinessDaysTrades(ctx context.Context) ([]models.Trade, error) {
	fromDate := s.getBusinessDaysAgo(7)
	toDate := time.Now()

	return s.GetTradesByDateRange(ctx, fromDate, toDate)
}

// getBusinessDaysAgo calculates the date N business days ago
func (s *TradeService) getBusinessDaysAgo(days int) time.Time {
	current := time.Now()
	businessDaysCount := 0

	for businessDaysCount < days {
		current = current.AddDate(0, 0, -1)
		// Skip weekends (Saturday = 6, Sunday = 0)
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			businessDaysCount++
		}
	}

	return current
}

// ValidateTicker checks if a ticker exists in the database
func (s *TradeService) ValidateTicker(ctx context.Context, ticker string) (bool, error) {
	tickers, err := s.tradeRepo.GetDistinctTickers(ctx)
	if err != nil {
		return false, err
	}

	for _, t := range tickers {
		if t == ticker {
			return true, nil
		}
	}

	return false, nil
}

// GetDistinctTickers returns all distinct tickers
func (s *TradeService) GetDistinctTickers(ctx context.Context) ([]string, error) {
	return s.tradeRepo.GetDistinctTickers(ctx)
}
