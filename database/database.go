package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Open database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}

	//Test Connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
		return nil, err
	}

	// Set Connection pool settings (optional, but recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connection established")
	return db, nil
}
