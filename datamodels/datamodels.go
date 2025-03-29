package datamodels

import (
	"time"
)

type User struct {
	CurrentUsername string `json:"current_username"`
	ID              int    `json:"id"`
}

type Profile struct {
	Username		string `json:"username"`
	ID				int    `json:"id"`
	RequestID		int    `json:"request_id"`
}

type Conversation struct {
	ConversationID string   `json:"conversation_id"`
	Participants   []string `json:"participants"`
	LastMessage    string   `json:"last_message,omitempty"` // omitempty allows the field to be optional
}

// Message represents a single message in a conversation.
type Message struct {
	ID             int            `json:"id"`
	Timestamp      time.Time      `json:"timestamp"`
	Preview        MessagePreview `json:"preview"`
	Comments       []Comment      `json:"comments"`
	MessageStatus  MessageStatus  `json:"message_status"`
	MessageContent MessageContent `json:"message_content"`
}

// MessagePreview represents the preview of a message.
type MessagePreview struct {
	Type         string `json:"type"` // "text" or "image"
	Content      string `json:"content"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
}

// Comment represents an emoji reaction to a message.
type Comment struct {
	Emoji     string    `json:"emoji"`
	UserID    int       `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

// MessageStatus represents the status of a message (received or sent).
type MessageStatus struct {
	Type           string `json:"type"` // "received" or "sent"
	SenderUsername string `json:"sender_username,omitempty"`
	Checkmarks     int    `json:"checkmarks,omitempty"`
}

// MessageContent represents the content of a message.
type MessageContent struct {
	Type     string `json:"type"` // "text" or "image"
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}