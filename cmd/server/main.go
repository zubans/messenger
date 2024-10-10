package main

import (
	"github.com/rs/cors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"video-conference/pkg/db"
	"video-conference/pkg/handlers"
)

func main() {
	// Подключение к базе данных
	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.CloseDB()

	// Конфигурация роутера
	r := mux.NewRouter()

	// Роуты для регистрации и входа
	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	// WebSocket
	r.Handle("/ws", http.HandlerFunc(handlers.SignalHandler))

	handler := cors.Default().Handler(r)

	// Запуск сервера
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
