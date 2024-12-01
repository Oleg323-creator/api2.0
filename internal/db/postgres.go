package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Migrator interface {
	Up() error
	Down() error
}

func NewPostgresDB(ctx context.Context, cfg ConnectionConfig) (*pgxpool.Pool, error) {

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host,
		cfg.Port, cfg.DBName, cfg.SSLMode)

	//PARCING CONFIG
	log.Printf("%s", connString)                 //(это для меня на будущее) can use easier way by pgxpool.Connect + defer dbPool.Close() without config parcing
	conf, err := pgxpool.ParseConfig(connString) // Using environment variables instead of a connection string.
	if err != nil {
		log.Fatalf("%s", err.Error())
		return nil, err
	}

	//CONNECTING CONFIG
	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		log.Fatalf("%s", err.Error())
		return nil, err
	}

	if err = getConnection(ctx, pool); err != nil {
		log.Fatalf("%s", err.Error())
		return nil, err
	}

	return pool, nil
}

// get connection from pool and release
func getConnection(ctx context.Context, pool *pgxpool.Pool) error {

	conn, err := pool.Acquire(ctx)

	defer conn.Release()

	if err != nil {
		log.Fatalf("%s", err.Error())
		return err
	}
	if err = conn.Ping(ctx); err != nil {
		log.Fatalf("%s", err.Error())
		return err
	}

	log.Println("Connected to database")
	return nil
}
