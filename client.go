package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func NewInlineKeyboardRow(buttons ...InlineKeyboardButton) []InlineKeyboardButton {
	var row []InlineKeyboardButton

	row = append(row, buttons...)

	return row
}

// SendMessage sends a message to a specific chat.
func (c *Client) SendMessage(chatID int64, text string) error {
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
func (c *Client) SendMarkup(chatID int64, text string, replyMarkup ReplyKeyboardMarkup) error {
	url := fmt.Sprintf("%s/sendMessage", c.BaseURL)

	payload := map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"reply_markup": replyMarkup,
	}

	return c.makeRequest(url, payload)
}

// SendInlineKeyboard sends a message with an inline keyboard to a specific chat.
func (c *Client) SendInlineKeyboard(chatID int64, text string, replyMarkup InlineKeyboardMarkup) error {
	url := fmt.Sprintf("%s/sendMessage", c.BaseURL)

	payload := map[string]interface{}{
		"chat_id":      chatID,
		"text":         text,
		"reply_markup": replyMarkup,
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

	// Log the response for debugging
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error response body: %s", respBody)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (m *Message) IsCommand() bool {
	for _, entity := range m.Entities {
		if entity.Type == "bot_command" {
			return true
		}
	}
	return false
}

// NewReplyKeyboard creates a new regular keyboard with sane defaults.
func NewReplyKeyboard(rows ...[]KeyboardButton) ReplyKeyboardMarkup {
	var keyboard [][]KeyboardButton

	keyboard = append(keyboard, rows...)

	return ReplyKeyboardMarkup{
		ResizeKeyboard: true,
		Keyboard:       keyboard,
	}
}

// NewKeyboardButtonRow creates a row of keyboard buttons.
func NewKeyboardButtonRow(buttons ...KeyboardButton) []KeyboardButton {
	var row []KeyboardButton

	row = append(row, buttons...)

	return row
}

// NewKeyboardButton creates a regular keyboard button.
func NewKeyboardButton(text string) KeyboardButton {
	return KeyboardButton{
		Text: text,
	}
}

// NewInlineKeyboardMarkup creates a new inline keyboard.
func NewInlineKeyboardMarkup(rows ...[]InlineKeyboardButton) InlineKeyboardMarkup {
	var keyboard [][]InlineKeyboardButton

	keyboard = append(keyboard, rows...)

	return InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}
}

// NewInlineKeyboardButtonData creates an inline keyboard button with text
// and data for a callback.
func NewInlineKeyboardButtonData(text, data string) InlineKeyboardButton {
	return InlineKeyboardButton{
		Text:         text,
		CallbackData: &data,
	}
}
