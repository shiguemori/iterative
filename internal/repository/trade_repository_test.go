package repository_test

import (
	"context"
	"testing"
	"time"

	"b3-challenge/internal/models"
	"b3-challenge/internal/repository"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// helper para criar db em memória
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	assert.NoError(t, err)

	err = db.AutoMigrate(&models.Trade{})
	assert.NoError(t, err)

	return db
}

func TestCreateTrades(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTradeRepository(db)

	trades := []models.Trade{
		{CodigoInstrumento: "ABC", PrecoNegocio: 10.5, QuantidadeNegociada: 100, DataNegocio: time.Now()},
		{CodigoInstrumento: "DEF", PrecoNegocio: 20.0, QuantidadeNegociada: 200, DataNegocio: time.Now()},
	}

	err := repo.CreateTrades(context.Background(), trades)
	assert.NoError(t, err)

	var count int64
	db.Model(&models.Trade{}).Count(&count)
	assert.Equal(t, int64(2), count)
}

func TestGetTradeAggregation(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTradeRepository(db)

	time1 := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
	time2 := time.Date(2023, 7, 2, 0, 0, 0, 0, time.UTC)
	ticker := "B3SA3"

	db.Create(&models.Trade{CodigoInstrumento: ticker, PrecoNegocio: 10.0, QuantidadeNegociada: 100, DataNegocio: time1})
	db.Create(&models.Trade{CodigoInstrumento: ticker, PrecoNegocio: 20.0, QuantidadeNegociada: 300, DataNegocio: time1})
	db.Create(&models.Trade{CodigoInstrumento: ticker, PrecoNegocio: 15.0, QuantidadeNegociada: 400, DataNegocio: time2})

	agg, err := repo.GetTradeAggregation(context.Background(), ticker, nil)
	assert.NoError(t, err)
	assert.Equal(t, 20.0, agg.MaxRangeValue)
	assert.Equal(t, int64(400), agg.MaxDailyVolume)
}

func TestGetTradesByDateRange(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTradeRepository(db)
	now := time.Now()

	db.Create(&models.Trade{CodigoInstrumento: "XPTO", DataNegocio: now.Add(-48 * time.Hour)})
	db.Create(&models.Trade{CodigoInstrumento: "XPTO", DataNegocio: now})

	from := now.Add(-72 * time.Hour)
	to := now.Add(1 * time.Hour)

	trades, err := repo.GetTradesByDateRange(context.Background(), from, to)
	assert.NoError(t, err)
	assert.Len(t, trades, 2)
}

func TestDeleteOldTrades(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTradeRepository(db)

	past := time.Now().Add(-48 * time.Hour)
	future := time.Now().Add(24 * time.Hour)

	db.Create(&models.Trade{CodigoInstrumento: "OLD", DataNegocio: past})
	db.Create(&models.Trade{CodigoInstrumento: "NEW", DataNegocio: future})

	err := repo.DeleteOldTrades(context.Background(), time.Now())
	assert.NoError(t, err)

	var trades []models.Trade
	db.Find(&trades)
	assert.Len(t, trades, 1)
	assert.Equal(t, "NEW", trades[0].CodigoInstrumento)
}

func TestGetDistinctTickers(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTradeRepository(db)

	db.Create(&models.Trade{CodigoInstrumento: "AAA"})
	db.Create(&models.Trade{CodigoInstrumento: "BBB"})
	db.Create(&models.Trade{CodigoInstrumento: "AAA"})

	tickers, err := repo.GetDistinctTickers(context.Background())
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"AAA", "BBB"}, tickers)
}
