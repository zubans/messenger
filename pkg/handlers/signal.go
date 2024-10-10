package handlers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"video-conference/pkg/db"
)

// Настройка upgrader для обновления HTTP соединений до WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Мапа для хранения подключенных клиентов
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte) // Канал для широковещательной передачи сообщений
var mutex = &sync.Mutex{}         // Мьютекс для защиты доступа к мапе клиентов

// Инициализация горутины для обработки сообщений
func init() {
	go handleMessages()
}

// Обработчик WebSocket соединений
func SignalHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	// Добавление нового клиента
	addClient(conn)

	// Чтение сообщений от клиента
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			removeClient(conn)
			break
		}

		// Отправка сообщения в канал broadcast
		broadcast <- message
	}
}

// Добавление клиента в мапу
func addClient(conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	clients[conn] = true
}

// Удаление клиента из мапы
func removeClient(conn *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(clients, conn)
}

// Обработка и трансляция сообщений всем подключенным клиентам
func handleMessages() {
	for {
		// Получаем следующее сообщение из канала broadcast
		message := <-broadcast

		// Широковещательная отправка сообщения всем подключенным клиентам
		mutex.Lock()
		for client := range clients {
			// Декодируем сообщение, чтобы извлечь токен
			var receivedMessage struct {
				Token   string `json:"token"`
				Message string `json:"message"`
			}

			if err := json.Unmarshal(message, &receivedMessage); err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			userID, err := validateAndExtractUserID(receivedMessage.Token)
			if err != nil {
				log.Printf("Invalid token: %v", err)
				continue
			}

			user, err := db.FindUserById(userID)
			if err != nil {
				log.Printf("Error finding user: %v", err)
				continue
			}

			fullMessage := struct {
				Username string `json:"username"`
				Message  string `json:"message"`
			}{
				Username: user.Username,
				Message:  receivedMessage.Message,
			}

			// Преобразуем и отправляем обогащенное сообщение
			finalMessage, err := json.Marshal(fullMessage)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			err = client.WriteMessage(websocket.TextMessage, finalMessage)
			if err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

// validateAndExtractUserID декодирует токен и извлекает user ID
func validateAndExtractUserID(tokenString string) (string, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	return claims.Subject, nil
}
