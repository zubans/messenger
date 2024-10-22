package db_test

import (
	"database/sql"
	"testing"
	"video-conference/pkg/db"

	"github.com/DATA-DOG/go-sqlmock"
	"video-conference/pkg/models"
)

func TestFindUserByUsername(t *testing.T) {
	// Создание Mock DB и структуры
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %s", err)
	}
	defer dbMock.Close()

	repo := &db.RepositoryImpl{DB: dbMock}

	t.Run("Find User Success", func(t *testing.T) {
		user := &models.User{
			ID:       "123",
			Username: "testuser",
			Email:    "testuser@example.com",
			Password: "hashedpassword",
		}

		// Предварительно настроим результат для запроса
		rows := sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(user.ID, user.Username, user.Email, user.Password)

		mock.ExpectQuery(`SELECT id, username, email, password FROM users WHERE username = \$1`).
			WithArgs("testuser").
			WillReturnRows(rows)

		// Вызываем тестируемую функцию
		result, err := repo.FindUserByUsername("testuser")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if result.Username != user.Username || result.ID != user.ID {
			t.Errorf("expected user %v, got %v", user, result)
		}
	})

	t.Run("Find User Not Found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, username, email, password FROM users WHERE username = \$1`).
			WithArgs("nonexistent").
			WillReturnError(sql.ErrNoRows)

		_, err := repo.FindUserByUsername("nonexistent")
		if err == nil {
			t.Fatalf("expected error, got none")
		}
	})
}

func TestSaveUser(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %s", err)
	}
	defer dbMock.Close()

	repo := &db.RepositoryImpl{DB: dbMock}

	t.Run("Save User Success", func(t *testing.T) {
		user := &models.User{
			ID:       "123",
			Username: "testuser",
			Email:    "testuser@example.com",
			Password: "hashedpassword",
		}

		mock.ExpectExec(`INSERT INTO users \(id, username, email, password\) VALUES \(\$1, \$2, \$3, \$4\)`).
			WithArgs(user.ID, user.Username, user.Email, user.Password).
			WillReturnResult(sqlmock.NewResult(1, 1))

		if err := repo.SaveUser(user); err != nil {
			t.Fatalf("unexpected error when saving user: %s", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("Save User Error", func(t *testing.T) {
		user := &models.User{
			ID:       "123",
			Username: "testuser",
			Email:    "testuser@example.com",
			Password: "hashedpassword",
		}

		mock.ExpectExec(`INSERT INTO users \(id, username, email, password\) VALUES \(\$1, \$2, \$3, \$4\)`).
			WithArgs(user.ID, user.Username, user.Email, user.Password).
			WillReturnError(sql.ErrConnDone) // Имитируем ошибку подключения

		err := repo.SaveUser(user)
		if err == nil {
			t.Fatalf("expected error when saving user, got none")
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
