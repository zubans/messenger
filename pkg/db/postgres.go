package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

//var DB *sql.DB

func (repo *RepositoryImpl) ConnectDB() error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"user",
		"password",
		"localhost",
		"5432",
		"videoconference",
	)
	var err error
	repo.DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := repo.DB.PingContext(ctx); err != nil {
		return err
	}

	log.Println("Connected to PostgreSQL database")
	return nil
}

func (repo *RepositoryImpl) CloseDB() {
	if repo.DB != nil {
		repo.DB.Close()
	}
}
