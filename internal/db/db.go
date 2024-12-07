package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DB CONFIG
type ConnectionConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// INIT CONNECTION TO DB
func NewDB(cfg ConnectionConfig) *sql.DB {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host,
		cfg.Port, cfg.DBName, cfg.SSLMode)
	log.Printf("Connecting to the database with connection string: %s", connString)

	// CONNECTING
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Error opening database connection:", err)

	}

	// CHECKING CONNECTION
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database:", err)
	}

	log.Println("Successfully connected to the database")
	return db
}

/*package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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

func NewDB(cfg ConnectionConfig) (*sql.DB, error) {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host,
		cfg.Port, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Проверка соединения
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connection to db", err)
	}

	// DRIVER INIT
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not initialize the postgres instance: %v", err)
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations", // Путь к папке с миграциями
		"postgres",            // Имя драйвера
		driver,                // Экземпляр драйвера
	)
	if err != nil {
		log.Fatalf("Migration failed: %v\n", err)
		return nil, err
	}

	// UP MIGRATIONS
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
		return nil, err
	}

	log.Println("Migrations applied successfully")
	return db, nil
}

/*
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
		log.Fatalf("error connection to db %s", err.Error())
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
*/
