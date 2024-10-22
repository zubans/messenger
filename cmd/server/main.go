package main

import (
	"github.com/rs/cors"
	"log"
	"net/http"
	"video-conference/pkg/auth"

	"github.com/gorilla/mux"
	"video-conference/pkg/db"
	"video-conference/pkg/handlers"
)

func main() {
	dbConn := &db.RepositoryImpl{}
	repo := handlers.Repos{}

	if err := dbConn.ConnectDB(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer dbConn.CloseDB()

	r := mux.NewRouter()

	// Маршруты, доступные без авторизации
	r.HandleFunc("/register", repo.Register).Methods("POST")
	r.HandleFunc("/login", repo.Login).Methods("POST")

	// Создаем подмаршруты, требующие аутентификации
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(auth.AuthMiddleware)

	// Защищенные маршруты
	protected.HandleFunc("/rooms", repo.GetRooms).Methods("GET")
	protected.HandleFunc("/rooms", repo.CreatePrivateRoom).Methods("POST")
	protected.HandleFunc("/users/search", repo.SearchUsers).Methods("GET")
	protected.Handle("/ws", http.HandlerFunc(repo.SignalHandler))

	// Настраиваем CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешаем все источники для демонстрации; укажите конкретные, как требуется
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(r)

	// Запуск сервера
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
