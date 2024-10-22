package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"video-conference/pkg/models"
	"video-conference/pkg/utils"
)

func (repos *Repos) Register(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user, err := models.NewUser(input.Username, input.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	if err := repos.UserRepo.SaveUser(user); err != nil {
		http.Error(w, "Failed to save user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User registered successfully")
}

func (repos *Repos) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	user, err := repos.UserRepo.FindUserByUsername(credentials.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	tokenString, err := utils.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Генерация нового UUID для токена
	tokenID, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, "Could not generate token ID", http.StatusInternalServerError)
		return
	}

	// Создание структуры токена
	token := &models.Token{
		ID:        tokenID.String(),
		Token:     tokenString,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}

	if err := repos.TokenRepo.SaveToken(token); err != nil {
		http.Error(w, "Failed to save token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
