package main

import (
	"context"
	"github.com/Oleg323-creator/api2.0/internal/db"
	"github.com/Oleg323-creator/api2.0/internal/runners"
	"github.com/Oleg323-creator/api2.0/internal/services"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file" // драйвер для работы с файлами миграций
	_ "github.com/golang-migrate/migrate/v4/source/file" // драйвер для файлов миграций
	_ "github.com/lib/pq"                                // драйвер для PostgreSQL
	_ "github.com/lib/pq"                                // драйвер для PostgreSQL
	"log"
	"sync"
)

const СoingeckoType = "Coingecko"
const CryptoCompType = "Crypto Compare"

func main() {
	// INIT CONFIG
	cfg := db.ConnectionConfig{
		Host:     "localhost",
		Port:     "5429",
		Username: "postgres",
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
	}

	// UP MIGRATIONS
	err := services.RunMigrations(cfg)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// DB INIT
	ctx := context.Background()
	dbConn, err := db.NewDB(ctx, "postgres", cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbConn.Close()

	runner, err := runners.NewRunner(CryptoCompType, 2, "BTC", "USDT", dbConn)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go runner.Run(ctx, &wg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal, stopping...")

	cancel()
	wg.Wait()

	log.Println("Application shut down gracefully")
}
