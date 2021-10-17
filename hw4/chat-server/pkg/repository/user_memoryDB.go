package repository

import (
	"chat/entities"
	"chat/memoryDB"
)

type UserMemoryDB struct {
	db *memoryDB.DB
}

func NewUserMemoryDB(db *memoryDB.DB) *UserMemoryDB {
	return &UserMemoryDB{db: db}
}

func (um *UserMemoryDB) SendMessageToUserByID(userID int, message entities.Message) error {
	return um.db.InsertMessageForUser(userID, message)
}

func (um *UserMemoryDB) GetUserMessages(userID int) ([]entities.Message, error) {
	return um.db.GetUserMessages(userID)
}


