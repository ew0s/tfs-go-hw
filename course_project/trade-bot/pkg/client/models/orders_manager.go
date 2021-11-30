package models

import (
	"fmt"
	"time"
)

type SendOrderInput struct {
	OrderType string `json:"order_type"`
	Symbol    string `json:"symbol"`
	Side      string `json:"side"`
	Size      uint   `json:"size"`
	JWTToken  string
}

type SendOrderResponse struct {
	ID                  string    `json:"id"`
	UserID              int       `json:"user_id"`
	ClientOrderID       string    `json:"client_order_id"`
	Type                string    `json:"type"`
	Symbol              string    `json:"symbol"`
	Quantity            int       `json:"quantity"`
	Side                string    `json:"side"`
	Filled              int       `json:"filled"`
	Timestamp           time.Time `json:"timestamp"`
	LastUpdateTimestamp time.Time `json:"last_update_timestamp"`
	Price               float64   `json:"price"`
	Message             string    `json:"message,omitempty"`
}

func (r *SendOrderResponse) String() string {
	if r.Message != "" {
		return fmt.Sprintf("Message: %s", r.Message)
	}

	return fmt.Sprintf(`
		order_id:   %s,
		type:       %s,
		symbol:     %s,
		quantity:   %d,
		side:       %s,
		filled:     %d,
		timestamp:  %s,
		price:      %f,
	`, r.ID, r.Type, r.Symbol, r.Quantity, r.Side, r.Filled, r.Timestamp, r.Price)
}

type StartTradingInput struct {
	OrderType        string  `json:"order_type"`
	Symbol           string  `json:"symbol"`
	Side             string  `json:"side"`
	Size             uint    `json:"size"`
	StopLossBorder   float64 `json:"stop_loss_border"`
	TakeProfitBorder float64 `json:"take_profit_border"`
	JWTToken         string
}

type StartTradingResponse struct {
	ID                  string    `json:"id"`
	UserID              int       `json:"user_id"`
	ClientOrderID       string    `json:"client_order_id"`
	Type                string    `json:"type"`
	Symbol              string    `json:"symbol"`
	Quantity            int       `json:"quantity"`
	Side                string    `json:"side"`
	Filled              int       `json:"filled"`
	Timestamp           time.Time `json:"timestamp"`
	LastUpdateTimestamp time.Time `json:"last_update_timestamp"`
	Price               float64   `json:"price"`
	Message             string    `json:"message,omitempty"`
}

func (r *StartTradingResponse) String() string {
	if r.Message != "" {
		return fmt.Sprintf("Message: %s", r.Message)
	}

	return fmt.Sprintf(`
		order_id:   %s,
		type:       %s,
		symbol:     %s,
		quantity:   %d,
		side:       %s,
		filled:     %d,
		timestamp:  %s,
		price:      %f,
	`, r.ID, r.Type, r.Symbol, r.Quantity, r.Side, r.Filled, r.Timestamp, r.Price)
}
