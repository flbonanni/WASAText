package database

import (
	"database/sql"
)

func (db *appdbimpl) CommentMessage(conversationId string, messageId string, emoji string, userID uint64) error {
	// Inserisce una nuova emoji reaction (comment) nella tabella comments
	res, err := db.c.Exec(
		`INSERT INTO comments (conversation_id, message_id, emoji, user_id, timestamp)
         VALUES (?, ?, ?, ?, ?)`,
		conversationId, messageId, emoji, userID, time.Now())
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// Assumi che ErrCommentNotCreated sia un errore definito altrove
		return ErrCommentNotCreated
	}
	return nil
}

func (db *appdbimpl) UncommentMessage(conversationId string, messageId string, userID uint64) error {
	// Elimina l'emoji reaction (comment) dalla tabella comments
	res, err := db.c.Exec(
		`DELETE FROM comments WHERE conversation_id = ? AND message_id = ? AND user_id = ?`,
		conversationId, messageId, userID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		// Assumi che ErrCommentDoesNotExist sia un errore definito altrove
		return ErrCommentDoesNotExist
	}
	return nil
}
