package services

import (
	"github.com/golang-migrate/migrate/v4"

	"fmt"
	"geckoapi1/internal/db"
	"log"
)

func RunMigrations(cfg db.ConnectionConfig) error {
	// CONNECTING TO MIGRASTIONS
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	// INITIALIZATION
	m, err := migrate.New(
		"file://./migrations",
		connString,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}

	// USING MIGRATIONS
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
