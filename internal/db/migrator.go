package db

import (
	"database/sql"
	_ "fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
)

// MIGRATIONS UP

func RunMigrations(db *sql.DB) error {
	// PSQL DRIVER INIT
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not initialize the postgres instance: %v", err)
		return err
	}

	// INIT MIGRATIONS
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations", // PATH TO DIRECTORY WITH MIGRATIONS
		"postgres",            // DB TYPE
		driver,                // DRIVER
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
		return err
	}

	// USE MIGRATIONS
	log.Println("Starting migrations...")
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}

/*package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // Это необходимо для работы с PostgreSQL
	"log"
)

func RunMigrations(cfg ConnectionConfig) error {
	// Строка подключения
	connString := fmt.Sprintf("postgresql://%s:%s@localhost:%s/%s?sslmode=%s", cfg.Username, cfg.Password, "5429", cfg.DBName, cfg.SSLMode)

	fmt.Println("Connection string:", connString)

	// OPEN DB
	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Failed to open the database: %v", err)
		return err
	}
	// CHECKING CONNECTION
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
		return err
	}

	// DRIVER INIT
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not initialize the postgres instance: %v", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations", // Путь к папке с миграциями
		"postgres",            // Имя драйвера
		driver,                // Экземпляр драйвера
	)
	if err != nil {
		log.Fatalf("Migration failed: %v\n", err)
		return err
	}

	// UP MIGRATIONS
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}
*/
