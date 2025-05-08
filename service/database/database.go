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
	"mime/multipart"
)

type User struct {
	Username        string `json:"username"` // new username vs current
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
	GetUserId(string) (User, error)

	GetConversations(string) ([]Conversation, error)
	GetConversation(string) (Conversation, error)

	UpdateGroupName(string, uint64, string) error
	UpdateGroupPhoto(string, uint64, multipart.File) error
	CreateGroup(uint64, string,  string, []string) (string, error)
	AddMemberToGroup(string, uint64, string) error
	RemoveMemberFromGroup(string, string) error

	CommentMessage(string, string, string, uint64) error
	UncommentMessage(string, string, uint64) error
	DeleteMessage(string, string, uint64) error
	SendMessage(string, Message) (Message, error)
	ForwardMessage(string, string, string, uint64) (Message, error)

	GetUserPicture(string) ([]byte, error)
	ChangeUserPhoto(User, Photo) error

	CreateUser(User) (User, error)
	SetUsername(User, string) (User, error)

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

    // mappa nome tabella → statement di creazione
    tables := map[string]string{
        "users": `
            CREATE TABLE IF NOT EXISTS users (
                id       INTEGER PRIMARY KEY,
                username TEXT    UNIQUE NOT NULL,
                photo    BLOB
            );
        `,
        "conversations": `
            CREATE TABLE IF NOT EXISTS conversations (
                conversation_id TEXT PRIMARY KEY,
                participants    TEXT NOT NULL,
                last_message    TEXT
            );
        `,
        "messages": `
            CREATE TABLE IF NOT EXISTS messages (
                id               INTEGER PRIMARY KEY,
                conversation_id  TEXT    NOT NULL,
                message_content  TEXT    NOT NULL,
                timestamp        DATETIME NOT NULL,
                sender_id        INTEGER NOT NULL,
                FOREIGN KEY(conversation_id) REFERENCES conversations(conversation_id),
                FOREIGN KEY(sender_id) REFERENCES users(id)
            );
        `,
        "comments": `
            CREATE TABLE IF NOT EXISTS comments (
                id              INTEGER PRIMARY KEY,
                conversation_id TEXT    NOT NULL,
                message_id      INTEGER NOT NULL,
                emoji           TEXT    NOT NULL,
                user_id         INTEGER NOT NULL,
                timestamp       DATETIME NOT NULL,
                FOREIGN KEY(conversation_id) REFERENCES conversations(conversation_id),
                FOREIGN KEY(message_id)      REFERENCES messages(id),
                FOREIGN KEY(user_id)         REFERENCES users(id)
            );
        `,
        "groups": `
            CREATE TABLE IF NOT EXISTS groups (
                group_id    TEXT    PRIMARY KEY,
                admin_id    INTEGER NOT NULL,
                group_name  TEXT    NOT NULL,
                description TEXT,
                members     TEXT    NOT NULL,  -- user1,user2,...
                photo       BLOB,
                FOREIGN KEY(admin_id) REFERENCES users(id)
            );
        `,
    }

    // esegue tutti i CREATE TABLE IF NOT EXISTS
    for tbl, stmt := range tables {
        if _, err := db.Exec(stmt); err != nil {
            return nil, fmt.Errorf("error creating table %q: %w", tbl, err)
        }
    }

    // alla fine, restituisci l’istanza pronta
    return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}