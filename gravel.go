package gravel

import "fmt"

type Gravel struct {
	Token   string
	BaseURL string
	Offset  int
	IsLog   bool
}

// NewGravel initializes a new Telegram bot client.
func NewGravel(token string, isLog bool) *Gravel {
	return &Gravel{
		Token:   token,
		BaseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
		IsLog:   isLog,
	}
}
