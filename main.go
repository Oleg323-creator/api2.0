package main

import (
	"context"
	"github.com/Oleg323-creator/api2.0/internal/runners"
	"os"
	"os/signal"
	"syscall"

	"github.com/Oleg323-creator/api2.0/internal/db"
	"github.com/Oleg323-creator/api2.0/internal/handlers"
	"log"
	"net/http"
	"sync"
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
	repo := db.NewRepository(dbConn)
	handler := handlers.NewHandler(repo)

	// RUN MIGRATIONS
	err := db.RunMigrations(dbConn)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	//SERVER
	mux := http.NewServeMux()
	mux.HandleFunc("/rates/list", handler.GetEndpoint)
	mux.HandleFunc("/rate/count", handler.PostEndpoint)
	mux.HandleFunc("/register", handler.SignUp)
	mux.HandleFunc("/login", handler.SignIn)

	handlerWithMiddleware := handlers.Middleware(mux)

	// INIT SERVER
	server := &http.Server{
		Addr:    ":8080",
		Handler: handlerWithMiddleware,
	}

	go func() {
		log.Println("Server is running on port 8080...")
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	runnerBtcCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "BTC", "USDT")
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerEthCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "ETH", "USDT")
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerBnbCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "BNB", "USDT")
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerUsdtBtcCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "USDT", "BTC")
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerUsdtEthCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "USDT", "ETH")
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	runnerUsdtBnbCryptoComp, err := runners.NewRunner(CryptoCompType, 1, "USDT", "BNB")
	if err != nil {
		log.Fatal("Failed to create runner:", err)
	}

	/*			//COINGECKO RUNNERS INIT
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
	*/
	runnerSlice := []*runners.Runner{runnerBtcCryptoComp, runnerUsdtBnbCryptoComp,
		runnerEthCryptoComp, runnerBnbCryptoComp, runnerUsdtBtcCryptoComp,
		runnerUsdtEthCryptoComp, /*runnerBnbCoingecko,, runnerUsdtEthCoingecko,
		runnerUsdtBnbCoingecko, runnerUsdtBtcСoingecko,
		runnerBtcCoingecko, runnerEthCoingecko*/
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	for _, runner := range runnerSlice {
		wg.Add(1)
		go func(r *runners.Runner) {
			defer wg.Done()
			r.Run(dbConn, ctx, &wg)
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
