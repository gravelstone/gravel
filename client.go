package gravel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gravelstone/gravel/logger"
)

// SendMessageToChannel sends a message to a Telegram channel.
func (c *Gravel) SendMessageToChannel(channelID string, text string) error {
	url := fmt.Sprintf("%s/sendMessage", c.BaseURL)

	payload := map[string]interface{}{
		"chat_id": channelID,
		"text":    text,
	}

	return c.makeRequest(url, payload)
}

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

func (c *Gravel) makeRequest(url string, payload map[string]interface{}) error {
	// Marshal the payload to JSON
	body, err := json.Marshal(payload)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to marshal payload: %v", err))
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	logger.Info(fmt.Sprintf("Sending request with payload: %s", string(body)))

	// Make the HTTP POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to send request: %v", err))
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read response body: %v", err))
		return fmt.Errorf("failed to read response body: %w", err)
	}

	logger.Info(fmt.Sprintf("Received response with status code: %d and body: %s", resp.StatusCode, string(respBody)))

	// Check if the status code indicates an error
	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Unexpected status code: %d with body: %s", resp.StatusCode, string(respBody)))
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	logger.Info("Request completed successfully")
	return nil
}

// GetUpdates fetches updates (messages) from the Telegram bot.
func (c *Gravel) GetUpdates() ([]Update, error) {
	// Build the request URL
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?offset=%d", c.Token, c.Offset)

	// Send the GET request
	resp, err := http.Get(url)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get updates: %v", err))
		return nil, fmt.Errorf("failed to send updates: %w", err)
	}
	defer resp.Body.Close()

	// Read and log the response body for debugging
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read response body: %v", err))
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response JSON
	var result struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result"`
	}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to decode response JSON: %v", err))
		return nil, fmt.Errorf("failed to decode response JSON: %w", err)
	}

	// Check the API response status
	if !result.Ok {
		logger.Error("Telegram returned an error")
		return nil, fmt.Errorf("telegram returned an error")
	}

	if len(result.Result) > 0 {
		c.Offset = result.Result[len(result.Result)-1].UpdateID + 1
	}

	logger.Info(fmt.Sprintf("Fetched %d updates, Received response: %s", len(result.Result), string(respBody)))
	return result.Result, nil
}

func (c *Gravel) GetUserInfo(chatID int64) (*User, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getChat?chat_id=%d", c.Token, chatID)

	resp, err := http.Get(url)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to send GET request: %v", err))
		return nil, fmt.Errorf("failed to send GET request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read response body: %v", err))
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the response JSON
	var result struct {
		Ok     bool `json:"ok"`
		Result Chat `json:"result"`
	}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to decode response JSON: %v", err))
		return nil, fmt.Errorf("failed to decode response JSON: %w", err)
	}

	if !result.Ok {
		logger.Error("Telegram returned an error")
		return nil, fmt.Errorf("telegram returned an error")
	}

	user := &User{
		ID:        result.Result.ID,
		FirstName: result.Result.FirstName,
		LastName:  result.Result.LastName,
		UserName:  result.Result.UserName,
	}

	logger.Info(fmt.Sprintf("Fetched user info: %+v", user))
	return user, nil
}
