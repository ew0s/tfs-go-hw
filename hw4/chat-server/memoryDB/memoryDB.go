package memoryDB

import (
	"chat/entities"
	"errors"
)

type userID int
type username string

type DB struct {
	users          map[userID]entities.User
	globalMessages []entities.Message
	usersUsernames map[username]userID
	usersMessages  map[userID][]entities.Message

	usersIdCount int
}

var (
	ErrUserAlreadyInDB = errors.New("user already in database")
	ErrNotFoundUser    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

func NewMemoryDB() (*DB, error) {
	db := DB{
		users:          map[userID]entities.User{},
		usersUsernames: map[username]userID{},
		usersMessages:  map[userID][]entities.Message{},
		usersIdCount:   1,
	}
	return &db, nil
}

func (db *DB) InsertUser(user entities.User) (int, error) {
	if _, ok := db.usersUsernames[username(user.Username)]; ok {
		return 0, ErrUserAlreadyInDB
	}

	user.ID = db.usersIdCount
	db.users[userID(user.ID)] = user
	db.insertUsersUsernames(username(user.Username), userID(user.ID))
	db.usersIdCount++
	return user.ID, nil
}

func (db *DB) InsertMessageToGlobalChat(usrID int, message entities.Message) error {
	if _, ok := db.users[userID(usrID)]; !ok {
		return ErrNotFoundUser
	}
	db.globalMessages = append(db.globalMessages, message)
	return nil
}

func (db *DB) GetMessagesFromGlobalChat() []entities.Message {
	return db.globalMessages
}

func (db *DB) InsertMessageForUser(usrID int, message entities.Message) error {
	_, ok := db.users[userID(usrID)]
	if !ok {
		return ErrNotFoundUser
	}
	db.usersMessages[userID(usrID)] = append(db.usersMessages[userID(usrID)], message)
	return nil
}

func (db *DB) GetUserMessages(usrID int) ([]entities.Message, error) {
	_, ok := db.users[userID(usrID)]
	if !ok {
		return nil, ErrNotFoundUser
	}
	return db.usersMessages[userID(usrID)], nil
}

func (db *DB) GetUser(usrName string, password string) (entities.User, error) {
	usrID, err := db.getUserIDByUsername(username(usrName))
	if err != nil {
		return entities.User{}, err
	}
	user, ok := db.users[usrID]
	if ok && user.Password == password {
		return user, nil
	} else if ok && user.Password != password {
		return entities.User{}, ErrInvalidPassword
	}
	return entities.User{}, ErrNotFoundUser
}

func (db *DB) getUserIDByUsername(usrName username) (userID, error) {
	usrID, ok := db.usersUsernames[usrName]
	if !ok {
		return 0, ErrNotFoundUser
	}
	return usrID, nil
}

func (db *DB) insertUsersUsernames(usrName username, usrID userID) {
	db.usersUsernames[usrName] = usrID
}
