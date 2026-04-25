package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "040290")
	dbname := getEnv("DB_NAME", "electronics_store")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

func RunMigrations() {
	dbURL := buildDatabaseURL()

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatal("Failed to create migrator:", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Migrations applied successfully")
}

func RollbackMigration() {
	dbURL := buildDatabaseURL()

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatal("Failed to create migrator:", err)
	}
	defer m.Close()

	if err := m.Steps(-1); err != nil {
		log.Fatal("Rollback failed:", err)
	}

	log.Println("Rolled back 1 migration step")
}

func buildDatabaseURL() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "040290")
	dbname := getEnv("DB_NAME", "electronics_store")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname,
	)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}