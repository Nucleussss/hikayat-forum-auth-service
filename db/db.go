package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func InitDB(dbString string) (*sql.DB, error) {
	const op = "db.InitDB"
	var db *sql.DB
	var err error
	var retriesMaxAttemp = 5

	log.Println("Initializing database connection")

	// connect to the database and handle any errors that occur
	for i := 0; i < retriesMaxAttemp; i++ {

		// open a new connection to the database
		db, err = sql.Open("postgres", dbString)
		if err != nil {
			log.Printf("%s Failed to connect to database (attemp %d): %v", op, i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		// check if the connection is valid by pinging the database
		if err = db.Ping(); err != nil {
			log.Printf("%s Failed to ping database (attemp %d): %v", op, i+1, err)

			// close the connection if it was opened successfully
			if closeErr := db.Close(); err != nil {
				log.Printf("%s Error closing database connection: %v", op, closeErr)
			}

			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Database connection established")
		return db, nil
	}

	return nil, fmt.Errorf("%s failed to connect to database after retries: %w", op, err)
}

func ConnectionString() string {
	DBHost := os.Getenv("DB_HOST")
	DBUser := os.Getenv("DB_USER")
	DBPass := os.Getenv("DB_PASSWORD")
	DBname := os.Getenv("DB_NAME")
	DBPort := os.Getenv("DB_PORT")
	DBsslMode := os.Getenv("DB_SSL_MODE")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", DBHost, DBUser, DBPass, DBname, DBPort, DBsslMode)
}
