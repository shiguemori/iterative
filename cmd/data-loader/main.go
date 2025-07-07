package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"b3-challenge/internal/config"
	"b3-challenge/internal/models"
	"b3-challenge/internal/repository"
	"b3-challenge/pkg/database"

	"github.com/joho/godotenv"
)

func main() {
	runtime.GOMAXPROCS(1)

	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg := config.Load()

	if err := database.Connect(&cfg.Database); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	tradeRepo := repository.NewTradeRepository(database.GetDB())

	dirPath := "./files"
	retryPath := filepath.Join(dirPath, "retry")
	donePath := filepath.Join(dirPath, "done")

	for _, path := range []string{retryPath, donePath} {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatalf("Error creating directory %s: %v", path, err)
		}
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Error reading directory %s: %v", dirPath, err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
			continue
		}
		fullPath := filepath.Join(dirPath, file.Name())
		processFile(ctx, fullPath, tradeRepo, retryPath, donePath)
	}

	log.Println("all files processed")
}

func processFile(ctx context.Context, path string, tradeRepo *repository.TradeRepository, retryPath, donePath string) {
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening file %s: %v", path, err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Error closing file %s: %v", path, err)
		} else {
			log.Printf("File %s closed successfully", path)
		}
	}(f)

	reader := csv.NewReader(f)
	reader.Comma = ';'
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading CSV file %s: %v", path, err)
		return
	}

	const batchSize = 1000
	const maxWorkers = 4

	var allRetryLines [][]string
	var mu sync.Mutex
	var wg sync.WaitGroup
	sema := make(chan struct{}, maxWorkers)

	for start := 1; start < len(records); start += batchSize {
		end := start + batchSize
		if end > len(records) {
			end = len(records)
		}
		batch := records[start:end]

		wg.Add(1)
		sema <- struct{}{}

		go func(batch [][]string) {
			defer wg.Done()
			defer func() { <-sema }()

			var localTrades []models.Trade
			var localRetry [][]string

			for _, row := range batch {
				dataReferencia, err1 := time.Parse("2006-01-02", row[0])
				precoNegocio, err2 := strconv.ParseFloat(strings.ReplaceAll(row[3], ",", "."), 64)
				quantidadeNegociada, err3 := strconv.ParseInt(row[4], 10, 64)
				dataNegocio, err4 := time.Parse("2006-01-02", row[8])

				if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
					localRetry = append(localRetry, row)
					continue
				}

				trade := models.Trade{
					DataReferencia:              dataReferencia,
					CodigoInstrumento:           row[1],
					AcaoAtualizacao:             row[2],
					PrecoNegocio:                precoNegocio,
					QuantidadeNegociada:         quantidadeNegociada,
					HoraFechamento:              row[5],
					CodigoIdentificadorNegocio:  row[6],
					TipoSessaoPregao:            row[7],
					DataNegocio:                 dataNegocio,
					CodigoParticipanteComprador: strings.TrimSpace(row[9]),
					CodigoParticipanteVendedor:  strings.TrimSpace(row[10]),
				}

				localTrades = append(localTrades, trade)
			}

			if len(localTrades) > 0 {
				if err = tradeRepo.CreateTrades(ctx, localTrades); err != nil {
					log.Printf("Error inserting batch trades from file %s: %v", path, err)
					localRetry = append(localRetry, batch...)
				} else {
					log.Printf("Successfully inserted %d trades from batch in file %s", len(localTrades), path)
				}
			}

			if len(localRetry) > 0 {
				mu.Lock()
				allRetryLines = append(allRetryLines, localRetry...)
				mu.Unlock()
			}
		}(batch)
	}

	wg.Wait()

	if len(allRetryLines) > 0 {
		baseName := filepath.Base(path)
		retryFile := filepath.Join(retryPath, baseName)
		saveRetryFile(retryFile, records[0], allRetryLines)
		log.Printf("Saved %d retry lines to %s", len(allRetryLines), retryFile)
	}

	dest := filepath.Join(donePath, filepath.Base(path))
	if err := os.Rename(path, dest); err != nil {
		log.Printf("Error moving file %s to done/: %v", path, err)
	}
}

func saveRetryFile(path string, header []string, rows [][]string) {
	f, err := os.Create(path)
	if err != nil {
		log.Printf("Error creating retry file %s: %v", path, err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("Error closing retry file %s: %v", path, err)
		} else {
			log.Printf("Retry file %s created successfully", path)
		}
	}(f)

	writer := csv.NewWriter(f)
	writer.Comma = ';'
	defer writer.Flush()

	err = writer.Write(header)
	if err != nil {
		return
	}
	err = writer.WriteAll(rows)
	if err != nil {
		return
	}
}
