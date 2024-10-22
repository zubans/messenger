package handlers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"video-conference/pkg/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (repo *Repos) SignalHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	roomID := r.URL.Query().Get("room_id")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		var receivedMessage struct {
			Token   string `json:"token"`
			Content string `json:"content"`
			RoomID  string `json:"room_id"`
		}

		if err := json.Unmarshal(msg, &receivedMessage); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Проверяем токен
		userID, err := validateAndExtractUserID(receivedMessage.Token)
		if err != nil {
			log.Printf("Invalid token: %v", err)
			continue
		}

		if err := repo.MessageRepo.SaveChatMessage(receivedMessage.Content, userID, roomID); err != nil {
			log.Printf("Error saving message to database: %v", err)
			continue
		}

		// Создаем структуру сообщения
		messageID := uuid.New().String()
		newMessage := &models.Message{
			ID:        messageID,
			Content:   receivedMessage.Content,
			UserID:    userID,
			RoomID:    receivedMessage.RoomID,
			CreatedAt: time.Now(),
		}

		// Сохранение сообщения в базу данных
		//if err := db.SaveMessage(newMessage); err != nil {
		//	log.Printf("Error saving message: %v", err)
		//	continue
		//}

		// Объединяем имя пользователя и передаем сообщение
		user, err := repo.UserRepo.FindUserById(userID)
		if err != nil {
			log.Printf("Error finding user: %v", err)
			continue
		}

		fullMessage := struct {
			Username string `json:"username"`
			Content  string `json:"content"`
			RoomID   string `json:"room_id"`
		}{
			Username: user.Username,
			Content:  newMessage.Content,
			RoomID:   newMessage.RoomID,
		}

		finalMessage, err := json.Marshal(fullMessage)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			continue
		}

		// Отправляем сообщение текущему соединению
		if err := conn.WriteMessage(websocket.TextMessage, finalMessage); err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}

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
