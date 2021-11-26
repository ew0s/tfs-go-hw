package models

type SendOrderInput struct {
	OrderType string `json:"order_type"`
	Symbol    string `json:"symbol"`
	Side      string `json:"side"`
	Size      int    `json:"size"`
	JwtToken  string
}

type SendOrderResponse struct {
	OrderID string `json:"order_id"`
	Message string `json:"message"`
}
