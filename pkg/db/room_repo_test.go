package db_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"
	"video-conference/pkg/db"
	"video-conference/pkg/models"
)

func TestGetRoomIDByName(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %s", err)
	}
	defer dbMock.Close()

	repo := &db.RepositoryImpl{DB: dbMock}

	t.Run("Room Exists", func(t *testing.T) {
		roomName := "Room 1"
		roomID := "12345"

		// Ожидание на успешное возвращение ID комнаты
		mock.ExpectQuery(`SELECT id FROM rooms WHERE name = \$1`).
			WithArgs(roomName).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(roomID))

		result, err := repo.GetRoomIDByName(roomName)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if result != roomID {
			t.Errorf("expected roomID %v, got %v", roomID, result)
		}
	})

	t.Run("Room Not Found", func(t *testing.T) {
		roomName := "NonExistentRoom"

		mock.ExpectQuery(`SELECT id FROM rooms WHERE name = \$1`).
			WithArgs(roomName).
			WillReturnError(sql.ErrNoRows)

		result, err := repo.GetRoomIDByName(roomName)
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		if result != "" {
			t.Errorf("expected empty string, got %v", result)
		}
	})
}

func TestSaveRoom(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %s", err)
	}
	defer dbMock.Close()

	repo := &db.RepositoryImpl{DB: dbMock}

	t.Run("Save Room Success", func(t *testing.T) {
		room := &models.Room{
			ID:        "12345",
			Name:      "Room 1",
			CreatedAt: time.Now(),
		}

		mock.ExpectExec(`INSERT INTO rooms \(id, name, created_at\) VALUES \(\$1, \$2, \$3\)`).
			WithArgs(room.ID, room.Name, room.CreatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SaveRoom(room)
		if err != nil {
			t.Fatalf("unexpected error when saving room: %s", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetRooms(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %s", err)
	}
	defer dbMock.Close()

	repo := &db.RepositoryImpl{DB: dbMock}

	t.Run("Get Rooms Success", func(t *testing.T) {
		mockRooms := sqlmock.NewRows([]string{"id", "name", "created_at"}).
			AddRow("12345", "Room 1", time.Now()).
			AddRow("67890", "Room 2", time.Now())

		mock.ExpectQuery(`SELECT id, name, created_at FROM rooms`).
			WillReturnRows(mockRooms)

		rooms, err := repo.GetRooms()
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if len(rooms) != 2 {
			t.Errorf("expected 2 rooms, got %d", len(rooms))
		}
	})
}
