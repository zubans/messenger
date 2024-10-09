package main

import (
	"log"
	"net/http"
	//"video-conference/pkg/db"
	"video-conference/pkg/signal"

	"github.com/gorilla/websocket"
)

const addr = ":8080"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	// Соединяемся с базой данных
	//if err := db.ConnectDB(); err != nil {
	//	log.Fatalf("Could not connect to the database: %v", err)
	//}
	//defer db.Conn.Close()

	// Обработчик для сигнализации WebRTC
	http.HandleFunc("/signal", signal.SignalHandler)

	// Запуск HTTP-сервера
	log.Printf("Server is listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

// Обработчик сигналов
func signalHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	//defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("Received: %s", message)
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
