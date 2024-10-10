package db

import (
	"log"
	"video-conference/pkg/models"
)

// SaveUser сохраняет нового пользователя в базу данных
func SaveUser(user *models.User) error {
	query := `INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)`
	_, err := DB.Exec(query, user.ID, user.Username, user.Email, user.Password)
	if err != nil {
		log.Printf("Error saving user: %v", err)
	}
	return err
}

// FindUserByUsername находит пользователя по имени пользователя
func FindUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE username = $1`

	row := DB.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return nil, err
	}

	return &user, nil
}

func FindUserById(id string) (*models.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE id = $1`

	row := DB.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return nil, err
	}

	return &user, nil
}
