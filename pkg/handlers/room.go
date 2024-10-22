package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
	"video-conference/pkg/auth"
	"video-conference/pkg/models"
)

// GetRooms возвращает список всех комнат
func (repo *Repos) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := repo.RoomRepo.GetRooms() // Вызов функции из db для получения списка комнат
	if err != nil {
		http.Error(w, "Could not retrieve rooms", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms) // Кодируем список комнат в JSON формат и отправляем клиенту
}

func (repo *Repos) CreatePrivateRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Извлечение currentUserID из контекста, добавленного middleware
	currentUserID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok {
		http.Error(w, "Could not retrieve user identity", http.StatusInternalServerError)
		return
	}

	// Создание уникального имени комнаты
	roomName := currentUserID + "-" + req.UserID
	roomID, err := repo.RoomRepo.GetRoomIDByName(roomName)
	if err == nil && roomID != "" {
		// Если комната уже существует, вернуть её ID
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"roomId": roomID})
		return
	}

	// Если комната не существует, создаем новую
	roomID = uuid.New().String()
	newRoom := &models.Room{
		ID:        roomID,
		Name:      roomName,
		CreatedAt: time.Now(),
	}

	if err := repo.RoomRepo.SaveRoom(newRoom); err != nil {
		http.Error(w, "Failed to create room", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"roomId": roomID})
}
