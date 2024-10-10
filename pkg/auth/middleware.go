package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

var jwtKey = []byte("your_secret_key") // Определите безопасный секретный ключ для вашей системы

// AuthMiddleware представляет middleware, который проверяет наличие и валидность JWT токена.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получение токена из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Извлечение токена, удаляя префикс "Bearer "
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

		// Разбор токена с использованием jwt-go
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		// Проверка валидности токена
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Если токен валидный, передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}
