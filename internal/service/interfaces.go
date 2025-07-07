package service

import (
	"context"
	"time"

	"b3-challenge/internal/models"
)

// TradeServiceInterface defines the interface for trade service operations
type TradeServiceInterface interface {
	GetTradeAggregation(ctx context.Context, ticker string) (*models.TradeAggregation, error)
	GetTradesByDateRange(ctx context.Context, fromDate, toDate time.Time) ([]models.Trade, error)
	GetLast7BusinessDaysTrades(ctx context.Context) ([]models.Trade, error)
	ValidateTicker(ctx context.Context, ticker string) (bool, error)
	GetDistinctTickers(ctx context.Context) ([]string, error)
}
