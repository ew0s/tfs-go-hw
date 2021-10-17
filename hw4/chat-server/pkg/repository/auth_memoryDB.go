package repository

import (
	"chat/entities"
	"chat/memoryDB"
)

type AuthMemoryDB struct {
	db *memoryDB.DB
}

func NewAuthMemory(db *memoryDB.DB) *AuthMemoryDB {
	return &AuthMemoryDB{db: db}
}

func (am *AuthMemoryDB) CreateUser(user entities.User) (int, error) {
	id, err := am.db.InsertUser(user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (am *AuthMemoryDB) GetUser(username string, password string) (entities.User, error) {
	var user entities.User
	user, err := am.db.GetUser(username, password)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}
