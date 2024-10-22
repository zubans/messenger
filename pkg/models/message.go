package models

import "time"

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	RoomID    string    `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
}
