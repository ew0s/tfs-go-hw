package entities

type Message struct {
	Content string `json:"content" binding:"required"`
	Sender  int
}
