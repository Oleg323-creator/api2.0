package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type ConnectionConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(ctx context.Context, cfg ConnectionConfig) (*pgxpool.Pool, error) {

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host,
		cfg.Port, cfg.DBName, cfg.SSLMode)

	//PARCING CONFIG
	log.Printf("%s", connString)                 //(это для меня на будущее) can use easier way by pgxpool.Connect + defer dbPool.Close() without config parcing
	conf, err := pgxpool.ParseConfig(connString) // Using environment variables instead of a connection string.
	if err != nil {
		log.Fatalf("Error %s", err.Error())
		return nil, err
	}

	//CONNECTING CONFIG
	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		log.Fatalf("%s", err.Error())
		return nil, err
	}

	err = getConnection(ctx, pool)
	if err != nil {
		log.Fatalf("%s", err.Error())
		return nil, err
	}

	return pool, nil
}

// get connection from pool and release
func getConnection(ctx context.Context, pool *pgxpool.Pool) error {

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("%s", err.Error())
		return err
	}

	defer conn.Release()

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalf("%s", err.Error())
		return err
	}

	log.Println("Connected to database")
	return nil
}

type WrapperDB struct {
	DBType string
	Pool   *pgxpool.Pool
	Ctx    context.Context
}

func NewDB(ctx context.Context, DBType string, cfg ConnectionConfig) (*WrapperDB, error) {
	var pool *pgxpool.Pool
	var err error
	switch DBType {
	case "postgres":
		pool, err = NewPostgresDB(ctx, cfg)
	default:
		pool, err = NewPostgresDB(ctx, cfg)
	}
	if err != nil {
		log.Fatalf("%s", err.Error())
		return nil, err
	}

	return &WrapperDB{
		DBType: DBType,
		Pool:   pool,
		Ctx:    ctx,
	}, nil
}

func (db *WrapperDB) Close() {
	db.Pool.Close()
}
