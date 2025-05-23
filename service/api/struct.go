package api

import (
	"github.com/flbonanni/WASAText/service/database"
) 
import "time"

type User struct {
	CurrentUsername string `json:"current_username"`
	ID              uint64 `json:"id"`
}

type Profile struct {
	Username		string `json:"username"`
	ID				uint64 `json:"id"`
	RequestID		uint64 `json:"request_id"`
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
	SenderID       string         `json:"sender_id"`
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

type Photo struct {
	Id            uint64 `json:"id"`
	UserId        uint64 `json:"userId"`
	File          []byte `json:"file"`
	Date          string `json:"date"`
}

func (u *User) FromDatabase(user database.User) {
	u.ID = user.ID
	u.CurrentUsername = user.CurrentUsername
}

func (u *User) ToDatabase() database.User {
	return database.User{
		ID:       u.ID,
		CurrentUsername: u.CurrentUsername,
	}
}

func (c *Conversation) ConvFromDatabase(conv Conversation) {
	c.ConversationID = conv.ConversationID
	c.Participants = conv.Participants
	c.LastMessage = conv.LastMessage
}
