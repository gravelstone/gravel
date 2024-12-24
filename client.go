package gravel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client represents the Telegram bot client.
type Client struct {
	Token   string
	BaseURL string
	Offset  int
}

// NewClient initializes a new Telegram bot client.
func NewClient(token string) *Client {
	return &Client{
		Token:   token,
		BaseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
	}
}

// SendMessage sends a message to a specific chat.
func (c *Client) SendMessage(chatID int64, text string) error {
	url := fmt.Sprintf("%s/sendMessage", c.BaseURL)

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	return c.makeRequest(url, payload)
}

// makeRequest is a helper function to send POST requests to the Telegram API.
func (c *Client) makeRequest(url string, payload map[string]interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
