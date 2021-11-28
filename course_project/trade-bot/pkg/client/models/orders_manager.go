package models

type SendOrderInput struct {
	OrderType string `json:"order_type"`
	Symbol    string `json:"symbol"`
	Side      string `json:"side"`
	Size      uint   `json:"size"`
	JWTToken  string
}

type SendOrderResponse struct {
	OrderID string `json:"order_id"`
	Message string `json:"message"`
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
	OrderID string `json:"order_id"`
}
