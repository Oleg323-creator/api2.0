package services

import (
	"database/sql"
	"fmt"
	"github.com/Oleg323-creator/api2.0/internal/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // Это необходимо для работы с PostgreSQL
	"log"
)

func RunMigrations(cfg db.ConnectionConfig) error {
	// Строка подключения
	connString := fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s?sslmode=%s", cfg.Username, cfg.Password, "5429", cfg.DBName, cfg.SSLMode)

	fmt.Println("Connection string:", connString)

	// Открытие подключения к базе данных
	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Failed to open the database: %v", err)
		return err
	}
	// Проверка, что подключение работает
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return err
	}

	// Инициализация драйвера
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not initialize the postgres instance: %v", err)
		return err
	}

	// Создание миграции с явно зарегистрированным драйвером
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations", // Путь к папке с миграциями
		"postgres",            // Имя драйвера
		driver,                // Экземпляр драйвера
	)
	if err != nil {
		log.Fatalf("Migration failed: %v\n", err)
		return err
	}

	// Применение миграций
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
