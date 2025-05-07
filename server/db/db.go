package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	log.Printf("Connecting to PostgreSQL on port 5433...")
	db, err := sql.Open("postgres", "postgresql://postgres:password@localhost:5433/go-chat?sslmode=disable")
	if err != nil {
		log.Printf("Error opening database: %v", err)
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
