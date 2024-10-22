package db_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"video-conference/pkg/db"
	"video-conference/pkg/models"
)

func TestSaveMessage(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %s", err)
	}
	defer dbMock.Close()

	repo := &db.RepositoryImpl{DB: dbMock}

	t.Run("Save Message Success", func(t *testing.T) {
		msg := &models.Message{
			ID:        "12345",
			Content:   "Hello world",
			UserID:    "user1",
			RoomID:    "room1",
			CreatedAt: time.Now(),
		}

		mock.ExpectExec(`INSERT INTO messages \(id, content, user_id, room_id, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
			WithArgs(msg.ID, msg.Content, msg.UserID, msg.RoomID, msg.CreatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SaveMessage(msg)
		if err != nil {
			t.Fatalf("unexpected error when saving message: %s", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestSaveChatMessage(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %s", err)
	}
	defer dbMock.Close()

	repo := &db.RepositoryImpl{DB: dbMock}

	t.Run("Save Chat Message Success", func(t *testing.T) {
		content := "Hello from chat"
		userID := "user2"
		roomID := "room2"

		mock.ExpectExec(`INSERT INTO messages \(id, content, user_id, room_id, created_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
			WithArgs(sqlmock.AnyArg(), content, userID, roomID, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.SaveChatMessage(content, userID, roomID)
		if err != nil {
			t.Fatalf("unexpected error when saving chat message: %s", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetMessagesForRoom(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock db: %s", err)
	}
	defer dbMock.Close()

	repo := &db.RepositoryImpl{DB: dbMock}

	t.Run("Get Messages for Room Success", func(t *testing.T) {
		roomID := "room1"
		mockMessages := sqlmock.NewRows([]string{"id", "content", "user_id", "room_id", "created_at"}).
			AddRow("1", "Message 1", "user1", roomID, time.Now()).
			AddRow("2", "Message 2", "user2", roomID, time.Now())

		mock.ExpectQuery(`SELECT id, content, user_id, room_id, created_at FROM messages WHERE room_id = \$1 ORDER BY created_at`).
			WithArgs(roomID).
			WillReturnRows(mockMessages)

		messages, err := repo.GetMessagesForRoom(roomID)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if len(messages) != 2 {
			t.Errorf("expected 2 messages, got %d", len(messages))
		}
	})
}
