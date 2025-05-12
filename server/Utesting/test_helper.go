package testing

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
)

var testDB *sql.DB

func SetupTestDB(t *testing.T) *sql.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file:", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Fatal("Error connecting to database:", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		t.Fatal("Error creating migration driver:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../db/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		t.Fatal("Error creating migrate instance:", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		t.Fatal("Error running migrations:", err)
	}

	return db
}

func CleanupTestDB(t *testing.T, db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing test database: %v", err)
		}
	}
}
