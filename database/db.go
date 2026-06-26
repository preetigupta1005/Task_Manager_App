package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file not loaded:", err)
	}

	conStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	db, err := sqlx.Open("postgres", conStr)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}

	DB = db
	log.Println("Connected to DB successfully")
	return nil
}

// close connection
func Close() {
	err := DB.Close()
	if err != nil {
		log.Fatal("DB not closed:", err)
	}
}
