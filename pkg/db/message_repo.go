package db

import (
	"github.com/google/uuid"
	"time"
	"video-conference/pkg/models"
)

// SaveMessage сохраняет сообщение в базу данных
func (repo *RepositoryImpl) SaveMessage(msg *models.Message) error {
	query := `INSERT INTO messages (id, content, user_id, room_id, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := repo.DB.Exec(query, msg.ID, msg.Content, msg.UserID, msg.RoomID, msg.CreatedAt)
	return err
}

func (repo *RepositoryImpl) SaveChatMessage(content string, userID string, roomID string) error {
	// Генерация нового UUID для сообщения
	messageID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	query := `INSERT INTO messages (id, content, user_id, room_id, created_at) VALUES ($1, $2, $3, $4, $5)`

	// Выполнение запрос в базу данных
	_, err = repo.DB.Exec(query, messageID.String(), content, userID, roomID, time.Now())
	return err
}

// GetMessagesForRoom получает сообщения для определенной комнаты
func (repo *RepositoryImpl) GetMessagesForRoom(roomID string) ([]models.Message, error) {
	query := `SELECT id, content, user_id, room_id, created_at FROM messages WHERE room_id = $1 ORDER BY created_at`
	rows, err := repo.DB.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.Content, &msg.UserID, &msg.RoomID, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
