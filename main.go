package main

import (
	"context"
	"geckoapi1/internal/db"
	"geckoapi1/internal/services"
	"github.com/joho/godotenv"
	"log"
	"os"
	"sync"
)

var (
	onceDB     sync.Once
	mgrWrapper *db.WrapperDB
	err        error
)

func main() {
	godotenv.Load()

	appCtx, cancel := context.WithCancel(context.Background())
	appWg := &sync.WaitGroup{} // new(sync.WaitGroup)???

	onceDB.Do(func() {
		mgrWrapper, err = db.NewDB(appCtx, "postgres", db.ConnectionConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
			SSLMode:  "disable",
		})
		if err != nil {
			log.Fatalf("error starting DB: %s", err.Error())
		}
		return
	})

	if err = services.RunMigrations(db.ConnectionConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  "disable",
	}); err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}

}
