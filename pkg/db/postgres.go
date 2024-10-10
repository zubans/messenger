package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"user",
		"password",
		"localhost",
		"5432",
		"videoconference",
	)
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		return err
	}

	log.Println("Connected to PostgreSQL database")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
