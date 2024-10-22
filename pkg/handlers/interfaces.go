package handlers

import (
	"video-conference/pkg/db"
)

type Repos struct {
	UserRepo    db.UserRepository
	TokenRepo   db.TokenRepository
	RoomRepo    db.RoomRepository
	MessageRepo db.MessageRepository
}
