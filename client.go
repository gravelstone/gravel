package gravel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Client represents the Telegram bot client.
func NewInlineKeyboardRow(buttons ...InlineKeyboardButton) []InlineKeyboardButton {
	var row []InlineKeyboardButton

	row = append(row, buttons...)

	return row
}

// SendMessage sends a message to a specific chat.
func (c *Gravel) SendMessage(chatID int64, text string) error {
	url := fmt.Sprintf("%s/sendMessage", c.BaseURL)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
		"reply_markup": map[string]interface{}{
			"remove_keyboard": true,
		},
	}

	return c.makeRequest(url, payload)
}

// SendMarkup sends a message with custom keyboard markup.
func (c *Gravel) SendMarkup(chatID int64, text string, replyMarkup ReplyKeyboardMarkup) error {
	url := fmt.Sprintf("%s/sendMessage", c.BaseURL)

	payload := map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"reply_markup": replyMarkup,
	}

	return c.makeRequest(url, payload)
}

// SendInlineKeyboard sends a message with an inline keyboard to a specific chat.
func (c *Gravel) SendInlineKeyboard(chatID int64, text string, replyMarkup InlineKeyboardMarkup) error {
	url := fmt.Sprintf("%s/sendMessage", c.BaseURL)

	payload := map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"reply_markup": replyMarkup,
	}

	return c.makeRequest(url, payload)
}

// makeRequest is a helper function to send POST requests to the Telegram API.
func (c *Gravel) makeRequest(url string, payload map[string]interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Log the response for debugging
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error response body: %s", respBody)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetUpdates fetches updates (messages) from the Telegram bot.
func (c *Gravel) GetUpdates() ([]Update, error) {
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

func (c *Gravel) GetUserInfo(chatID int64) (*User, error) {
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
