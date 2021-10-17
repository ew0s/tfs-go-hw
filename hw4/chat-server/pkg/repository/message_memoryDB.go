package repository

import (
	"chat/entities"
	"chat/memoryDB"
)

type MessageMemoryDB struct {
	db *memoryDB.DB
}

func NewMessageMemoryDB(db *memoryDB.DB) *MessageMemoryDB {
	return &MessageMemoryDB{db: db}
}

func (mm *MessageMemoryDB) Create(usrID int, message entities.Message) error {
	return mm.db.InsertMessageToGlobalChat(usrID, message)
}

func (mm *MessageMemoryDB) GetGlobalMessages() []entities.Message {
	return mm.db.GetMessagesFromGlobalChat()
}
