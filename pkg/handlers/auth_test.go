package handlers_test

//
//import (
//	"bytes"
//	"encoding/json"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//	"video-conference/pkg/handlers"
//
//	"video-conference/pkg/db"
//	"video-conference/pkg/models"
//
//	"github.com/golang/mock/gomock"
//	"golang.org/x/crypto/bcrypt"
//)
//
//func TestRegister(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	mockUserRepo := db.NewMockUserRepository(ctrl)
//
//	defer ctrl.Finish()
//
//	// Настроим mock для сохранения пользователя
//	mockUserRepo.EXPECT().SaveUser(gomock.Any()).Return(nil)
//
//	t.Run("Success Register", func(t *testing.T) {
//		payload := models.User{
//			Username: "testuser",
//			Email:    "testuser@example.com",
//			Password: "password123",
//		}
//		body, _ := json.Marshal(payload)
//
//		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
//		if err != nil {
//			t.Fatal(err)
//		}
//		rr := httptest.NewRecorder()
//
//		handler := http.HandlerFunc(handlers.Register)
//		handler.ServeHTTP(rr, req)
//
//		if status := rr.Code; status != http.StatusCreated {
//			t.Errorf("expected %v, got %v", http.StatusCreated, status)
//		}
//
//		expected := `"User registered successfully"`
//		if rr.Body.String() != expected {
//			t.Errorf("expected %v, got %v", expected, rr.Body.String())
//		}
//	})
//
//	// Другие случаи для Register тут
//}
//
//func TestLogin(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	mockUserRepo := db.NewMockUserRepository(ctrl)
//
//	defer ctrl.Finish()
//
//	// Подделываем успешный возврат пользователя
//	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
//	mockUserRepo.EXPECT().FindUserByUsername("testuser").Return(&models.User{
//		ID:       "1234",
//		Username: "testuser",
//		Password: string(hashedPassword),
//	}, nil).AnyTimes()
//
//	t.Run("Success Login", func(t *testing.T) {
//		payload := map[string]string{
//			"username": "testuser",
//			"password": "correct-password",
//		}
//		body, _ := json.Marshal(payload)
//
//		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
//		if err != nil {
//			t.Fatal(err)
//		}
//		rr := httptest.NewRecorder()
//
//		handler := http.HandlerFunc(handlers.Login)
//		handler.ServeHTTP(rr, req)
//
//		if status := rr.Code; status != http.StatusOK {
//			t.Errorf("expected %v, got %v", http.StatusOK, status)
//		}
//
//		var response map[string]string
//		json.Unmarshal(rr.Body.Bytes(), &response)
//
//		if _, ok := response["token"]; !ok {
//			t.Errorf("expected a token in response")
//		}
//	})
//
//	// Другие случаи для Login тут
//}
