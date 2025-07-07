package service

import (
	"context"
	"log"
	"math/rand"
	"time"

	"b3-challenge/internal/models"
	"b3-challenge/internal/repository"
)

// MockDataService generates mock data for demonstration purposes
type MockDataService struct {
	tradeRepo *repository.TradeRepository
}

// NewMockDataService creates a new mock data service
func NewMockDataService(tradeRepo *repository.TradeRepository) *MockDataService {
	return &MockDataService{
		tradeRepo: tradeRepo,
	}
}

// GenerateMockData creates mock trade data for the last 7 business days
func (s *MockDataService) GenerateMockData(ctx context.Context) error {
	log.Println("Generating mock data for demonstration...")

	tickers := []string{
		"PETR4", "VALE3", "ITUB4", "BBDC4", "ABEV3", "MGLU3", "WEGE3",
		"RENT3", "LREN3", "SUZB3", "RAIL3", "USIM5", "CSNA3", "GOAU4",
		"BBAS3", "SANB11", "BPAC11", "ITSA4", "CIEL3", "GGBR4",
	}

	// Base prices for each ticker (realistic values)
	basePrices := map[string]float64{
		"PETR4":  38.50,
		"VALE3":  65.20,
		"ITUB4":  32.15,
		"BBDC4":  14.80,
		"ABEV3":  11.25,
		"MGLU3":  8.90,
		"WEGE3":  45.30,
		"RENT3":  58.75,
		"LREN3":  18.40,
		"SUZB3":  12.60,
		"RAIL3":  19.85,
		"USIM5":  7.45,
		"CSNA3":  13.20,
		"GOAU4":  9.80,
		"BBAS3":  28.90,
		"SANB11": 35.60,
		"BPAC11": 22.40,
		"ITSA4":  9.15,
		"CIEL3":  4.25,
		"GGBR4":  18.70,
	}

	// Generate data for last 7 business days
	businessDays := s.getLastBusinessDays(7)

	var allTrades []models.Trade

	for _, ticker := range tickers {
		basePrice := basePrices[ticker]
		currentPrice := basePrice

		for _, day := range businessDays {
			// Generate multiple trades per day for each ticker
			tradesPerDay := rand.Intn(50) + 10 // 10-60 trades per day

			for i := 0; i < tradesPerDay; i++ {
				// Simulate price variation (±5%)
				variation := (rand.Float64() - 0.5) * 0.1 // ±5%
				tradePrice := currentPrice * (1 + variation)

				// Ensure price doesn't go negative
				if tradePrice < 0.01 {
					tradePrice = 0.01
				}

				// Generate random quantity (100-10000 shares)
				quantity := int64(rand.Intn(9900) + 100)

				// Generate random time during trading hours (9:00-17:30)
				hour := rand.Intn(8) + 9 // 9-16
				minute := rand.Intn(60)
				second := rand.Intn(60)

				tradeTime := time.Date(day.Year(), day.Month(), day.Day(), hour, minute, second, 0, time.UTC)
				closingTime := time.Date(0, 1, 1, hour, minute, second, 0, time.UTC).Format("15:04:05")

				trade := models.Trade{
					DataReferencia:      tradeTime,
					CodigoInstrumento:   ticker,
					QuantidadeNegociada: quantity,
					PrecoNegocio:        tradePrice,
					HoraFechamento:      closingTime,
				}

				allTrades = append(allTrades, trade)
			}

			// Update current price for next day (small drift)
			drift := (rand.Float64() - 0.5) * 0.02 // ±1% daily drift
			currentPrice = currentPrice * (1 + drift)
		}
	}

	// Save all trades to database
	if len(allTrades) > 0 {
		if err := s.tradeRepo.CreateTrades(ctx, allTrades); err != nil {
			return err
		}
		log.Printf("Generated and saved %d mock trades", len(allTrades))
	}

	return nil
}

// getLastBusinessDays returns the last N business days (excluding weekends)
func (s *MockDataService) getLastBusinessDays(days int) []time.Time {
	var businessDays []time.Time
	current := time.Now()

	for len(businessDays) < days {
		// Skip weekends (Saturday = 6, Sunday = 0)
		if current.Weekday() != time.Saturday && current.Weekday() != time.Sunday {
			businessDays = append(businessDays, current)
		}
		current = current.AddDate(0, 0, -1)
	}

	// Reverse to get chronological order
	for i := 0; i < len(businessDays)/2; i++ {
		j := len(businessDays) - 1 - i
		businessDays[i], businessDays[j] = businessDays[j], businessDays[i]
	}

	return businessDays
}
