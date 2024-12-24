package gravel

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetUpdates fetches updates (messages) from the Telegram bot.
func (c *Client) GetUpdates() ([]Update, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d", c.Token, c.Offset)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if !result.Ok {
		return nil, fmt.Errorf("telegram API returned an error")
	}

	if len(result.Result) > 0 {
		c.Offset = result.Result[len(result.Result)-1].UpdateID + 1
	}

	return result.Result, nil
}
