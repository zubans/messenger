package utils

import (
	"github.com/golang-jwt/jwt/v4"
	"time"

	"video-conference/pkg/models"
)

var jwtKey = []byte("your_secret_key") // Можете заменить на более безопасный ключ

// GenerateJWT генерирует новый JWT токен для пользователя
func GenerateJWT(user *models.User) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Subject:   user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
