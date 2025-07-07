package repository

import (
	"context"
	"time"

	"b3-challenge/internal/models"
)

// TradeRepositoryInterface defines the interface for trade repository operations
type TradeRepositoryInterface interface {
	CreateTrades(ctx context.Context, trades []models.Trade) error
	GetTradeAggregation(ctx context.Context, ticker string, fromDate *time.Time) (*models.TradeAggregation, error)
	GetTradesByDateRange(ctx context.Context, fromDate, toDate time.Time) ([]models.Trade, error)
	DeleteOldTrades(ctx context.Context, beforeDate time.Time) error
	GetDistinctTickers(ctx context.Context) ([]string, error)
}
