package repository

import (
	"errors"
	"fmt"
	"trade-bot/configs"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // postgres driver import
)

var (
	ErrNewPostgresDB = errors.New("new postgres db")
	ErrPingDB        = errors.New("ping db")
)

func NewPostgresDB(cfg configs.DatabaseConfiguration) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ErrNewPostgresDB, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %s: %w", ErrNewPostgresDB, ErrPingDB, err)
	}

	return db, nil
}
