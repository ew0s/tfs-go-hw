package postgresRepo

import (
	"trade-bot/internal/pkg/models"

	"github.com/jmoiron/sqlx"
)

type KrakenOrdersManagerPostgres struct {
	db *sqlx.DB
}

func NewKrakenOrdersManagerPostgres(db *sqlx.DB) *KrakenOrdersManagerPostgres {
	return &KrakenOrdersManagerPostgres{db: db}
}

const createOrderQuery = `
	INSERT INTO orders(order_id, user_id, cli_order_id, type, symbol, quantity, side, filled,
	                  timestamp, last_update_timestamp, limit_price)
	VALUES(:order_id, :user_id, :cli_order_id, :type, :symbol, :quantity, :side, :filled,
	                  :timestamp, :last_update_timestamp, :limit_price)`

const createUsersOrdersQuery = `
	INSERT INTO users_orders(user_id, order_id) VALUES ($1, $2)
`

func (k *KrakenOrdersManagerPostgres) CreateOrder(userID int, order models.Order) error {
	_, err := k.db.NamedExec(createOrderQuery, order)
	if err != nil {
		return err
	}

	_, err = k.db.Exec(createUsersOrdersQuery, userID, order.ID)
	return err
}
