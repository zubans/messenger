package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

var Conn *pgx.Conn

func ConnectDB() error {
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		"user",
		"password",
		"localhost", // IP-адрес или имя хоста контейнера postgres
		"5432",
		"videoconference",
	))
	if err != nil {
		return err
	}

	Conn = conn
	log.Println("Connected to PostgreSQL database!")
	return nil
}
