package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` // Хранить хэшированный пароль
}

// NewUser создает нового пользователя с уникальным UUID.
func NewUser(username, email, password string) (*User, error) {
	// Генерация нового UUID
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       id.String(),
		Username: username,
		Email:    email,
		Password: password, // Важно хэшировать пароль перед сохранением
	}, nil
}
