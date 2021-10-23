package memoryDB

import (
	"chat/entities"
	"errors"
	"sync"
)

type userID int
type username string

type DB struct {
	usersMutex          sync.RWMutex
	globalMessagesMutex sync.RWMutex
	usersUsernamesMutex sync.RWMutex
	usersMessagesMutex  sync.RWMutex
	users               map[userID]entities.User
	globalMessages      []entities.Message
	usersUsernames      map[username]userID
	usersMessages       map[userID][]entities.Message

	usersIDCount int
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
		usersIDCount:   1,
	}
	return &db, nil
}

func (db *DB) InsertUser(user entities.User) (int, error) {
	db.usersMutex.Lock()
	db.usersUsernamesMutex.Lock()
	defer db.usersMutex.Unlock()
	defer db.usersUsernamesMutex.Unlock()
	return db.insertUser(user)
}

func (db *DB) InsertMessageToGlobalChat(usrID int, message entities.Message) error {
	db.usersMutex.RLock()
	db.globalMessagesMutex.Lock()
	defer db.usersMutex.RUnlock()
	defer db.globalMessagesMutex.Unlock()
	return db.insertMessageToGlobalChat(usrID, message)
}

func (db *DB) InsertMessageForUser(usrID int, message entities.Message) error {
	db.usersMutex.RLock()
	db.usersMessagesMutex.Lock()
	defer db.usersMutex.RUnlock()
	defer db.usersMessagesMutex.Unlock()
	return db.insertMessageForUser(usrID, message)
}

func (db *DB) GetMessagesFromGlobalChat() []entities.Message {
	db.globalMessagesMutex.RLock()
	defer db.globalMessagesMutex.RUnlock()
	return db.getMessagesFromGlobalChat()
}

func (db *DB) GetUserMessages(usrID int) ([]entities.Message, error) {
	db.usersMutex.RLock()
	db.usersMessagesMutex.RLock()
	defer db.usersMutex.RUnlock()
	defer db.usersMessagesMutex.RUnlock()
	return db.getUserMessages(usrID)
}

func (db *DB) GetUser(usrName string, password string) (entities.User, error) {
	db.usersMutex.RLock()
	db.usersUsernamesMutex.RLock()
	defer db.usersMutex.RUnlock()
	defer db.usersUsernamesMutex.RUnlock()
	return db.getUser(usrName, password)
}

func (db *DB) insertUser(user entities.User) (int, error) {
	if _, ok := db.usersUsernames[username(user.Username)]; ok {
		return 0, ErrUserAlreadyInDB
	}

	user.ID = db.usersIDCount
	db.users[userID(user.ID)] = user
	db.usersUsernames[username(user.Username)] = userID(user.ID)
	db.usersIDCount++
	return user.ID, nil
}

func (db *DB) insertMessageToGlobalChat(usrID int, message entities.Message) error {
	if _, ok := db.users[userID(usrID)]; !ok {
		return ErrNotFoundUser
	}
	db.globalMessages = append(db.globalMessages, message)
	return nil
}

func (db *DB) getMessagesFromGlobalChat() []entities.Message {
	return db.globalMessages
}

func (db *DB) insertMessageForUser(usrID int, message entities.Message) error {
	_, ok := db.users[userID(usrID)]
	if !ok {
		return ErrNotFoundUser
	}
	db.usersMessages[userID(usrID)] = append(db.usersMessages[userID(usrID)], message)
	return nil
}

func (db *DB) getUserMessages(usrID int) ([]entities.Message, error) {
	_, ok := db.users[userID(usrID)]
	if !ok {
		return nil, ErrNotFoundUser
	}
	return db.usersMessages[userID(usrID)], nil
}

func (db *DB) getUser(usrName string, password string) (entities.User, error) {
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
