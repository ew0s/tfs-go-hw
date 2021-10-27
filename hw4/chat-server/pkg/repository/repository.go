package repository

import (
	"chat/entities"
	"chat/memoryDB"
)

type Authorization interface {
	CreateUser(user entities.User) (int, error)
	GetUser(username string, password string) (entities.User, error)
}

type Message interface {
	Create(usrID int, message entities.Message) error
	GetGlobalMessages() []entities.Message
}

type User interface {
	SendMessageToUserByID(userID int, message entities.Message) error
	GetUserMessages(userID int) ([]entities.Message, error)
}

type Repository struct {
	Authorization
	Message
	User
}

func NewRepository(db *memoryDB.DB) *Repository {
	return &Repository{
		Authorization: NewAuthMemory(db),
		Message:       NewMessageMemoryDB(db),
		User:          NewUserMemoryDB(db),
	}
}
