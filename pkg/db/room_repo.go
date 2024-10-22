package db

import (
	"database/sql"
	"video-conference/pkg/models"
)

func (repo *RepositoryImpl) GetRoomIDByName(name string) (string, error) {
	query := `SELECT id FROM rooms WHERE name = $1`
	row := repo.DB.QueryRow(query, name)

	var roomID string
	if err := row.Scan(&roomID); err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Комната не найдена, возвращаем пустую строку
		}
		return "", err // Возвращаем ошибку, если она произошла
	}

	return roomID, nil
}

func (repo *RepositoryImpl) SaveRoom(room *models.Room) error {
	query := `INSERT INTO rooms (id, name, created_at) VALUES ($1, $2, $3)`
	_, err := repo.DB.Exec(query, room.ID, room.Name, room.CreatedAt)
	return err
}

// GetRooms извлекает список всех комнат
func (repo *RepositoryImpl) GetRooms() ([]models.Room, error) {
	query := `SELECT id, name, created_at FROM rooms`

	rows, err := repo.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		if err := rows.Scan(&room.ID, &room.Name, &room.CreatedAt); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}
