package db_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"video-conference/pkg/db"
	"video-conference/pkg/models"
)

func TestSaveToken(t *testing.T) {
	// Создание Mock DB и структуры
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %s", err)
	}
	defer dbMock.Close()

	repo := &db.RepositoryImpl{DB: dbMock}

	t.Run("Save Token Success", func(t *testing.T) {
		token := &models.Token{
			ID:        "12345",
			Token:     "sometoken",
			UserID:    "67890",
			CreatedAt: time.Now(),
		}

		// Ожидание на успешное выполнение Exec
		mock.ExpectExec(`INSERT INTO tokens \(id, token, user_id, created_at\) VALUES \(\$1, \$2, \$3, \$4\)`).
			WithArgs(token.ID, token.Token, token.UserID, token.CreatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SaveToken(token)
		if err != nil {
			t.Fatalf("unexpected error when saving token: %s", err)
		}

		// Проверка, все ли ожидания выполнены
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
