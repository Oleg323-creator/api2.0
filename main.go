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

	//RUNNERS INIT

	runnerBtcCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "BTC", "USDT", dbConn)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerEthCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "ETH", "USDT", dbConn)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerBnbCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "BNB", "USDT", dbConn)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}
	/*
		runnerBtcCoingecko, err := runners.NewRunner(СoingeckoType, 1, "BTC", "USDT", dbConn)
		if err != nil {
			log.Fatal("Failed to create runner:", err)
		}

		runnerEthCoingecko, err := runners.NewRunner(СoingeckoType, 1, "ETH", "USDT", dbConn)
		if err != nil {
			log.Fatal("Failed to create runner:", err)
		}

		runnerBnbCoingecko, err := runners.NewRunner(СoingeckoType, 1, "BNB", "USDT", dbConn)
		if err != nil {
			log.Fatal("Failed to create runner:", err)
		}

		//OPOSITE RUNNERS INIT

		runnerUsdtBtcСoingecko, err := runners.NewRunner(СoingeckoType, 1, "USDT", "BTC", dbConn)
		if err != nil {
			log.Fatal("Failed to create runner:", err)
		}

		runnerUsdtEthCoingecko, err := runners.NewRunner(СoingeckoType, 1, "USDT", "ETH", dbConn)
		if err != nil {
			log.Fatal("Failed to create runner:", err)
		}

		runnerUsdtBnbCoingecko, err := runners.NewRunner(СoingeckoType, 1, "USDT", "BNB", dbConn)
		if err != nil {
			log.Fatal("Failed to create runner:", err)
		}
	*/
	runnerUsdtBtcCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "USDT", "BTC", dbConn)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerUsdtEthCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "USDT", "ETH", dbConn)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerUsdtBnbCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "USDT", "BNB", dbConn)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerSlice := []*runners.Runner{runnerBtcCryptoComp, runnerUsdtBnbCryptoComp,
		runnerEthCryptoComp, runnerBnbCryptoComp,
		runnerUsdtBtcCryptoComp, runnerUsdtEthCryptoComp,
		/*runnerUsdtEthCoingecko, runnerUsdtBnbCoingecko, runnerBtcCoingecko, runnerEthCoingecko,
		runnerBnbCoingecko, runnerUsdtBtcСoingecko*/}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for _, runner := range runnerSlice {
		wg.Add(1)
		go runner.Run(ctx, &wg)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal, stopping...")

	cancel()
	wg.Wait()

	log.Println("Application shut down gracefully")
}
