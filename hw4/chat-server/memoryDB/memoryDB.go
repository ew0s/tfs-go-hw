package memoryDB

import (
	"chat/entities"
	"errors"
	"fmt"
	"sync"
)

type userID int
type username string

type dbUsers struct {
	sync.RWMutex
	data map[userID]entities.User
}

type dbGlobalMessages struct {
	sync.RWMutex
	data []entities.Message
}

type dbUsersUsernames struct {
	sync.RWMutex
	data map[username]userID
}

type dbUsersMessages struct {
	sync.RWMutex
	data map[userID][]entities.Message
}

type DB struct {
	users          dbUsers
	globalMessages dbGlobalMessages
	usersUsernames dbUsersUsernames
	usersMessages  dbUsersMessages

	usersIDCount int
}

var (
	ErrInsertUser                = errors.New("insert user")
	ErrInsertMessageToGlobalChat = errors.New("inset message to global chat")
	ErrUserAlreadyInDB           = errors.New("user already in database")
	ErrInsertMessageForUser      = errors.New("insert message for user")
	ErrGetUserMessages           = errors.New("get user messages")
	ErrGetUser                   = errors.New("get user")
	ErrNotFoundUser              = errors.New("user not found")
	ErrUsernameNotFound          = errors.New("username not found")
	ErrUsernameAlreadyInDB       = errors.New("username already in db")
	ErrInvalidPassword           = errors.New("invalid password")
)

func NewMemoryDB() (*DB, error) {
	db := DB{
		users: dbUsers{
			data: map[userID]entities.User{},
		},
		usersUsernames: dbUsersUsernames{
			data: map[username]userID{},
		},
		usersMessages: dbUsersMessages{
			data: map[userID][]entities.Message{},
		},
		usersIDCount: 1,
	}
	return &db, nil
}

func (db *DB) InsertUser(user entities.User) (int, error) {
	user.ID = db.usersIDCount
	if err := db.usersUsernames.insert(username(user.Username), userID(user.ID)); err != nil {
		return 0, fmt.Errorf("%s: %w", ErrInsertUser, err)
	}
	if err := db.users.insert(user); err != nil {
		return 0, fmt.Errorf("%s: %w", ErrInsertUser, err)
	}

	db.usersIDCount++
	return user.ID, nil
}

func (db *DB) InsertMessageToGlobalChat(usrID int, message entities.Message) error {
	if _, err := db.users.get(userID(usrID)); err != nil {
		return fmt.Errorf("%s: %w", ErrInsertMessageToGlobalChat, err)
	}
	db.globalMessages.insert(message)
	return nil
}

func (db *DB) GetMessagesFromGlobalChat() []entities.Message {
	return db.globalMessages.get()
}

func (db *DB) InsertMessageForUser(usrID int, message entities.Message) error {
	if _, err := db.users.get(userID(usrID)); err != nil {
		return fmt.Errorf("%s: %w", ErrInsertMessageForUser, err)
	}
	db.usersMessages.insert(userID(usrID), message)
	return nil
}

func (db *DB) GetUserMessages(usrID int) ([]entities.Message, error) {
	if _, err := db.users.get(userID(usrID)); err != nil {
		return nil, fmt.Errorf("%s: %w", ErrGetUserMessages, err)
	}
	return db.usersMessages.get(userID(usrID)), nil
}

func (db *DB) GetUser(usrName string, password string) (entities.User, error) {
	usrID, err := db.usersUsernames.get(username(usrName))
	if err != nil {
		return entities.User{}, fmt.Errorf("%s: %w", ErrGetUser, err)
	}
	user, err := db.users.get(usrID)
	if err != nil {
		return entities.User{}, fmt.Errorf("%s: %w", ErrGetUser, err)
	}
	if user.Password == password {
		return user, nil
	}
	return entities.User{}, fmt.Errorf("%s: %s", ErrGetUser, ErrInvalidPassword)
}

func (u *dbUsers) insert(user entities.User) error {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.data[userID(user.ID)]; ok {
		return ErrUserAlreadyInDB
	}
	u.data[userID(user.ID)] = user
	return nil
}

func (u *dbUsers) get(usrID userID) (entities.User, error) {
	u.RLock()
	defer u.RUnlock()
	usr, ok := u.data[usrID]
	if !ok {
		return entities.User{}, ErrNotFoundUser
	}
	return usr, nil
}

func (m *dbGlobalMessages) insert(message entities.Message) {
	m.Lock()
	defer m.Unlock()
	m.data = append(m.data, message)
}

func (m *dbGlobalMessages) get() []entities.Message {
	m.RLock()
	defer m.RUnlock()
	messagesCopy := make([]entities.Message, len(m.data))
	copy(messagesCopy, m.data)
	return messagesCopy
}

func (u *dbUsersUsernames) insert(usrName username, usrID userID) error {
	u.Lock()
	defer u.Unlock()
	if _, ok := u.data[usrName]; ok {
		return ErrUsernameAlreadyInDB
	}
	u.data[usrName] = usrID
	return nil
}

func (u *dbUsersUsernames) get(usrName username) (userID, error) {
	u.RLock()
	defer u.RUnlock()
	usrID, ok := u.data[usrName]
	if !ok {
		return 0, ErrUsernameNotFound
	}
	return usrID, nil
}

func (m *dbUsersMessages) insert(usrID userID, message entities.Message) {
	m.Lock()
	defer m.Unlock()
	m.data[usrID] = append(m.data[usrID], message)
}

func (m *dbUsersMessages) get(usrID userID) []entities.Message {
	m.RLock()
	defer m.RUnlock()
	if val, ok := m.data[usrID]; ok && val != nil {
		messagesCopy := make([]entities.Message, len(val))
		copy(messagesCopy, val)
		return messagesCopy
	}
	return []entities.Message{}
}
