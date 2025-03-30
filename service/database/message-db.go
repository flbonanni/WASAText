package database

import (
	"database/sql"
	"github.com/flbonanni/WASAText/datamodels"
	"time"
	"fmt"
)

func (db *appdbimpl) SendMessage(conversationId string, m datamodels.Message) (datamodels.Message, error) {
	// Inserisce il messaggio nel database
	res, err := db.c.Exec(
		`INSERT INTO messages (conversation_id, message_content, timestamp, sender_id)
         VALUES (?, ?, ?, ?)`,
		conversationId, m.MessageContent, m.Timestamp)
	if err != nil {
		return m, err
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return m, err
	}
	m.ID = int(lastInsertID)
	return m, nil
}

func (db *appdbimpl) ForwardMessage(messageId string, targetConversationId string, recipientUsername string, senderID uint64) (datamodels.Message, error) {
	var orig datamodels.Message

	// Recupera il messaggio originale dalla tabella messages
	err := db.c.QueryRow(
		`SELECT id, message_content, timestamp FROM messages WHERE id = ?`,
		messageId).Scan(&orig.ID, &orig.MessageContent, &orig.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return datamodels.Message{}, fmt.Errorf("message not found")
		}
		return datamodels.Message{}, err
	}

	// Se si desidera modificare il contenuto in forward (ad esempio, aggiungere un prefisso) si può fare qui.
	forwardedContent := orig.MessageContent

	// Inserisce il messaggio inoltrato nella conversazione di destinazione
	now := time.Now()
	res, err := db.c.Exec(
		`INSERT INTO messages (conversation_id, message_content, timestamp, sender_id)
         VALUES (?, ?, ?, ?)`,
		targetConversationId, forwardedContent, now, senderID)
	if err != nil {
		return datamodels.Message{}, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return datamodels.Message{}, err
	}

	forwardedMessage := datamodels.Message{
		ID:             int(lastInsertID),
		MessageContent: forwardedContent,
		Timestamp:      now,
	}

	return forwardedMessage, nil
}

func (db *appdbimpl) DeleteMessage(conversationId string, messageId string, senderID uint64) error {
	// Elimina il messaggio verificando che appartenga alla conversazione e sia stato inviato dall'utente
	res, err := db.c.Exec(
		`DELETE FROM messages WHERE id = ? AND conversation_id = ? AND sender_id = ?`,
		messageId, conversationId, senderID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("message not found") // oppure un errore di autorizzazione, a seconda della logica
	}
	return nil
}
