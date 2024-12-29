package gravel

import (
	"fmt"

	"github.com/gravelstone/gravel/logger"
)

type Gravel struct {
	Token   string
	BaseURL string
	Offset  int
}

// NewGravel initializes a new Telegram bot client.
func NewGravel(token string, isLog bool) *Gravel {
	if isLog {
		logger.Config(true)
	}

	return &Gravel{
		Token:   token,
		BaseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
	}
}
