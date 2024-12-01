package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type ConnectionConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
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
	switch db.DBType {
	case "postgres":
		db.Pool.Close()
	default:
		db.Pool.Close()
	}
}

// Функция для записи данных в базу данных
func saveDataToDB(data Data) error {
	query := "INSERT INTO rates (value, timestamp) VALUES ($1, $2)"
	_, err := dbPool.Exec(context.Background(), query, data.Value, data.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert data: %v", err)
	}
	log.Println("Data saved to DB:", data)
	return nil
}

// Периодическая функция, генерирующая данные и записывающая их в базу
func generateAndSaveData() {
	data := Data{
		Value:     "RandomValue", // Здесь ваши реальные данные
		Timestamp: time.Now(),
	}

	err := saveDataToDB(data)
	if err != nil {
		log.Printf("Error saving data: %v", err)
	}
}
