package main

import (
	"context"
	"github.com/Oleg323-creator/api2.0/internal/db"
	"github.com/Oleg323-creator/api2.0/internal/db/rep"
	"github.com/Oleg323-creator/api2.0/internal/handlers"
	"github.com/Oleg323-creator/api2.0/internal/runners"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const СoingeckoType = "Coingecko"
const CryptoCompType = "Crypto_Compare"

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

	// DB CONNECT
	dbConn := db.NewDB(cfg)

	//CONNECT DB FOR MAKING REQUSTS
	repo := rep.NewRepository(dbConn)
	handler := handlers.NewHandler(repo)

	// RUN MIGRATIONS
	err := db.RunMigrations(dbConn)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/rates/list", handlers.AuthenticationMiddleware(), handler.GetEndpoint)
	router.POST("/rate/count", handlers.AuthenticationMiddleware(), handler.PostEndpoint)
	router.POST("/register", handler.SignUp)
	router.POST("/login", handler.SignIn)

	go func() {
		err := router.Run(":8080")
		if err != nil {
			log.Fatal("Failed to start Gin server:", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Initialize cryptocurrency "runners"
	curr := []string{"BTC", "ETH", "BNB"}
	for i := range curr {
		_, err = runners.NewRunner(CryptoCompType, 1, curr[i], "USDT", repo)
		if err != nil {
			log.Fatal("Failed to create runner:", err)
			return
		}
	}

	for i := range curr {
		runner, err := runners.NewRunner(CryptoCompType, 1, "USDT", curr[i], repo)
		if err != nil {
			log.Fatal("Failed to create runner:", err)
			return
		}
		wg.Add(1)
		go func(r *runners.Runner) {
			defer wg.Done()
			r.Run(ctx, &wg)
		}(runner)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal, stopping...")

	log.Println("Application shut down gracefully")

	cancel()
	wg.Wait()
}

/*
		runnerBtcCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "BTC", "USDT", repo)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerEthCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "ETH", "USDT", repo)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerBnbCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "BNB", "USDT", repo)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerUsdtBtcCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "USDT", "BTC", repo)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerUsdtEthCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "USDT", "ETH", repo)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerUsdtBnbCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "USDT", "BNB", repo)
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

		//COINGECKO RUNNERS INIT
	runnerBnbCoingecko, err := runners.NewRunner(СoingeckoType, 1, "BNB", "USDT")
	if err != nil {
		log.Fatal("Failed to create runner:", err)
		}

	runnerBtcCoingecko, err := runners.NewRunner(СoingeckoType, 1, "BTC", "USDT")
	if err != nil {
		log.Fatal("Failed to create runner:", err)
		}

	runnerEthCoingecko, err := runners.NewRunner(СoingeckoType, 1, "ETH", "USDT")
	if err != nil {
		log.Fatal("Failed to create runner:", err)
		}

	runnerUsdtBtcСoingecko, err := runners.NewRunner(СoingeckoType, 1, "USDT", "BTC")
	if err != nil {
		log.Fatal("Failed to create runner:", err)
		}

	runnerUsdtBnbCoingecko, err := runners.NewRunner(СoingeckoType, 1, "USDT", "BNB")
	if err != nil {
		log.Fatal("Failed to create runner:", err)
		}

	runnerUsdtEthCoingecko, err := runners.NewRunner(СoingeckoType, 1, "USDT", "ETH")
	if err != nil {
		log.Fatal("Failed to create runner:", err)
		}

	runnerSlice := []*runners.Runner{runnerBtcCryptoComp, runnerUsdtBnbCryptoComp,
		runnerEthCryptoComp, runnerBnbCryptoComp, runnerUsdtBtcCryptoComp,
		runnerUsdtEthCryptoComp, /*runnerBnbCoingecko,, runnerUsdtEthCoingecko,
		runnerUsdtBnbCoingecko, runnerUsdtBtcСoingecko,
		runnerBtcCoingecko, runnerEthCoingecko*/

/*	for _, runner := range runnerSlice {
		wg.Add(1)
		go func(r *runners.Runner) {
			defer wg.Done()
			r.Run(ctx, &wg)
		}(runner)
	}
*/
