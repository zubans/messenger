package signal

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrader для обновления HTTP-соединений до WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Хранение активных соединений
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)
var mutex = &sync.Mutex{}

// Инициируем горутину для обработки входящих сообщений
func init() {
	go handleMessages()
}

// SignalHandler обрабатывает новые подключения WebSocket
func SignalHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	addClient(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			removeClient(conn)
			break
		}
		broadcast <- message
	}
}

// Добавление нового клиента
func addClient(conn *websocket.Conn) {
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()
}

// Удаление клиента
func removeClient(conn *websocket.Conn) {
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
}

// handleMessages отправляет входящие сообщения всем подключенным клиентам
func handleMessages() {
	for {
		// Получаем следующее сообщение из канала broadcast
		message := <-broadcast

		// Отправка сообщений всем клиентам
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
