package database

import (
	"database/sql"
    "encoding/json"
	"time"
)

func (db *appdbimpl) SendMessage(conversationId string, m Message) (Message, error) {
    // Serializziamo MessageContent in JSON
    contentBytes, err := json.Marshal(m.MessageContent)
    if err != nil {
        return m, err
    }

    // Inseriamo il messaggio; qui assumo che la tabella messages abbia le colonne
    // (conversation_id, message_content, timestamp), senza sender_id
    res, err := db.c.Exec(
        `INSERT INTO messages (conversation_id, message_content, timestamp, sender_id)
        VALUES (?, ?, ?, ?)`,
        conversationId,
        string(contentBytes),
        m.Timestamp,
        m.SenderID,
     )

    lastInsertID, err := res.LastInsertId()
    if err != nil {
        return m, err
    }
    m.ID = int(lastInsertID)
    return m, nil
}

func (db *appdbimpl) ForwardMessage(
    messageId string,
    targetConversationId string,
    recipientUsername string,
    senderID string, // ora stringa, coerente con SenderID nel modello
) (Message, error) {
    var orig Message

    // 1) Recupero del messaggio originale in una stringa
    var contentStr string
    err := db.c.QueryRow(
        `SELECT id, message_content, timestamp, sender_id 
           FROM messages 
          WHERE id = ?`,
        messageId,
    ).Scan(&orig.ID, &contentStr, &orig.Timestamp, &orig.SenderID)
    if err != nil {
        if err == sql.ErrNoRows {
            return orig, ErrMessageDoesNotExist
        }
        return orig, err
    }

    // 2) Unmarshal del JSON in orig.MessageContent
    if err := json.Unmarshal([]byte(contentStr), &orig.MessageContent); err != nil {
        return orig, err
    }

    // 3) (Opzionale) Modifica del contenuto per il forward
    forwardedContent := orig.MessageContent

    // 4) Serializzo di nuovo il MessageContent in JSON
    forwardBytes, err := json.Marshal(forwardedContent)
    if err != nil {
        return orig, err
    }

    // 5) Inserimento nella conversazione di destinazione
    now := time.Now()
    res, err := db.c.Exec(
        `INSERT INTO messages (conversation_id, message_content, timestamp, sender_id)
         VALUES (?, ?, ?, ?)`,
        targetConversationId,
        string(forwardBytes),
        now,
        senderID,
    )
    if err != nil {
        return orig, err
    }

    newID, err := res.LastInsertId()
    if err != nil {
        return orig, err
    }

    // 6) Costruzione del messaggio inoltrato da restituire
    forwardedMsg := Message{
        ID:             int(newID),
        Timestamp:      now,
        SenderID:       senderID,
        MessageContent: forwardedContent,
        // Preview, Comments, MessageStatus li puoi lasciare vuoti o popolare se ti servono
    }

    return forwardedMsg, nil
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
		return ErrMessageDoesNotExist
	}
	return nil
}
