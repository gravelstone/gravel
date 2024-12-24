package telegram

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

func (c *Client) GetUserInfo(chatID int64) (*User, error) {
	// Construct the URL for the Telegram Bot API
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getChat?chat_id=%d", c.Token, chatID)

	// Make the HTTP request to get user info
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Define a struct to hold the response
	var result struct {
		Ok     bool `json:"ok"`
		Result Chat `json:"result"`
	}

	// Decode the JSON response into the result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	// Check if the response was successful
	if !result.Ok {
		return nil, fmt.Errorf("telegram API returned an error")
	}

	// Convert the chat info into user info (You can adjust fields as needed)
	user := &User{
		ID:        result.Result.ID,
		FirstName: result.Result.FirstName,
		LastName:  result.Result.LastName,
		UserName:  result.Result.UserName,
	}

	return user, nil
}
