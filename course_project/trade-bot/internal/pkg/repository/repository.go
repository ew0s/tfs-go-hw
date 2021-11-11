package repository

import (
	"trade-bot/internal/pkg/models"
	"trade-bot/internal/pkg/repository/postgresRepo"
	"trade-bot/internal/pkg/repository/redisRepo"
	"trade-bot/pkg/utils"

	"github.com/go-redis/redis/v8"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(models.User) (int, error)
	GetUser(username string) (models.User, error)
	GetUserAPIKeys(userID int) (string, string, error)
}

type JWT interface {
	CreateJWT(userID int, td utils.TokenDetails) (string, error)
	GetJWTUserID(ad utils.AccessDetails) (int, error)
	DeleteJWT(ad utils.AccessDetails) error
}

type Repository struct {
	Authorization
	JWT
}

func NewRepository(db *sqlx.DB, jwtDB *redis.Client) *Repository {
	return &Repository{
		Authorization: postgresRepo.NewAuthPostgres(db),
		JWT:           redisRepo.NewJWTRedis(jwtDB),
	}
}
