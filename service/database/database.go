/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	CurrentUsername string `json:"current_username"`
	ID              uint64 `json:"id"`
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

type Photo struct {
	Id            uint64 `json:"id"`
	UserId        uint64 `json:"userId"`
	File          []byte `json:"file"`
	Date          string `json:"date"`
}

var ErrUserDoesNotExist = errors.New("User does not exist")
var ErrPhotoDoesNotExist = errors.New("Photo does not exist")
var ErrBanDoesNotExist = errors.New("Ban does not exist")
var ErrFollowDoesNotExist = errors.New("Follow does not exist")
var ErrCommentDoesNotExist = errors.New("Comment does not exist")
var ErrLikeDoesNotExist = errors.New("Like does not exist")
var ErrMessageDoesNotExist = errors.New("Message does not exist")

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(string) error
	CheckUserById(User) (User, error)
	CommentMessage(string, string, string, uint64) error
	GetConversations(string) ([]Conversation, error)
	UncommentMessage(string, string, uint64) error
	GetConversation(string) (Conversation, error)
	GetUserId(string) (User, error)
	UpdateGroupName(string, uint64, string) error
	UpdateGroupPhoto(string, uint64, string) error
	CreateGroup(uint64, string,  string, []string) (string, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}