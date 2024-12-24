package gravel

// Update represents a Telegram update object.
type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message"`
}

// Message represents a Telegram message object.
type Message struct {
	MessageID int    `json:"message_id"`
	Chat      *Chat  `json:"chat"`
	Text      string `json:"text"`
}

// Chat represents a Telegram chat object.
type Chat struct {
	ID int64 `json:"id"`
}
