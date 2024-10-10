package db

import (
	"video-conference/pkg/models"
)

// SaveToken сохраняет новый токен в базе данных
func SaveToken(token *models.Token) error {
	query := `INSERT INTO tokens (id, token, user_id, created_at) VALUES ($1, $2, $3, $4)`
	_, err := DB.Exec(query, token.ID, token.Token, token.UserID, token.CreatedAt)
	return err
}
