package db

import (
	"database/sql"
	"video-conference/pkg/models"
)

type RepositoryImpl struct {
	DB *sql.DB
}

// UserRepository определяет операции с пользователями
type UserRepository interface {
	FindUserByUsername(username string) (*models.User, error)
	SaveUser(user *models.User) error
	FindUsersByUsername(name string) ([]models.User, error)
	FindUserById(id string) (*models.User, error)
}

// TokenRepository определяет операции с токенами
type TokenRepository interface {
	SaveToken(token *models.Token) error
}

type RoomRepository interface {
	GetRoomIDByName(name string) (string, error)
	SaveRoom(room *models.Room) error
	GetRooms() ([]models.Room, error)
}

type MessageRepository interface {
	SaveMessage(message *models.Message) error
	SaveChatMessage(content string, userID string, roomID string) error
	GetMessagesForRoom(roomID string) ([]models.Message, error)
}
