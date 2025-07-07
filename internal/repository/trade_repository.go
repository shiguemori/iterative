package repository

import (
	"context"
	"fmt"
	"time"

	"b3-challenge/internal/models"

	"gorm.io/gorm"
)

// TradeRepository handles database operations for trades
type TradeRepository struct {
	db *gorm.DB
}

// NewTradeRepository creates a new trade repository
func NewTradeRepository(db *gorm.DB) *TradeRepository {
	return &TradeRepository{db: db}
}

// CreateTrades inserts multiple trades in batch
func (r *TradeRepository) CreateTrades(ctx context.Context, trades []models.Trade) error {
	if len(trades) == 0 {
		return nil
	}

	// Use batch insert for better performance
	batchSize := 1000
	if err := r.db.WithContext(ctx).CreateInBatches(trades, batchSize).Error; err != nil {
		return fmt.Errorf("failed to create trades batch: %w", err)
	}

	return nil
}

// GetTradeAggregation returns aggregated trade data for a specific ticker
func (r *TradeRepository) GetTradeAggregation(ctx context.Context, ticker string, fromDate *time.Time) (*models.TradeAggregation, error) {
	var result models.TradeAggregation

	query := r.db.WithContext(ctx).Model(&models.Trade{}).
		Where("codigo_instrumento = ?", ticker)

	if fromDate != nil {
		query = query.Where("data_negocio >= ?", *fromDate)
	}

	// Get max trade price (max_range_value)
	var maxPrice float64
	if err := query.Select("MAX(preco_negocio)").Scan(&maxPrice).Error; err != nil {
		return nil, fmt.Errorf("failed to get max trade price: %w", err)
	}

	// Get max daily volume (max_daily_volume)
	var maxDailyVolume int64
	subQuery := r.db.Model(&models.Trade{}).
		Select("SUM(quantidade_negociada) as daily_volume").
		Where("codigo_instrumento = ?", ticker).
		Group("data_negocio")

	if fromDate != nil {
		subQuery = subQuery.Where("data_negocio >= ?", *fromDate)
	}

	if err := r.db.WithContext(ctx).Table("(?) as daily_volumes", subQuery).
		Select("MAX(daily_volume)").Scan(&maxDailyVolume).Error; err != nil {
		return nil, fmt.Errorf("failed to get max daily volume: %w", err)
	}

	result = models.TradeAggregation{
		Ticker:         ticker,
		MaxRangeValue:  maxPrice,
		MaxDailyVolume: maxDailyVolume,
	}

	return &result, nil
}

// GetTradesByDateRange returns trades within a date range
func (r *TradeRepository) GetTradesByDateRange(ctx context.Context, fromDate, toDate time.Time) ([]models.Trade, error) {
	var trades []models.Trade

	err := r.db.WithContext(ctx).
		Where("data_negocio BETWEEN ? AND ?", fromDate, toDate).
		Order("data_negocio DESC, codigo_instrumento").
		Find(&trades).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get trades by date range: %w", err)
	}

	return trades, nil
}

// DeleteOldTrades removes trades older than specified date
func (r *TradeRepository) DeleteOldTrades(ctx context.Context, beforeDate time.Time) error {
	result := r.db.WithContext(ctx).
		Where("data_negocio < ?", beforeDate).
		Delete(&models.Trade{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete old trades: %w", result.Error)
	}

	return nil
}

// GetDistinctTickers returns all unique tickers in the database
func (r *TradeRepository) GetDistinctTickers(ctx context.Context) ([]string, error) {
	var tickers []string

	err := r.db.WithContext(ctx).
		Model(&models.Trade{}).
		Distinct("codigo_instrumento").
		Pluck("codigo_instrumento", &tickers).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get distinct tickers: %w", err)
	}

	return tickers, nil
}
