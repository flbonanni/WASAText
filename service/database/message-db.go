package database

import (
	"database/sql"
    "encoding/json"
	"time"
	"strconv"
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
    senderID uint64,                      // ora uint64, come nell’interfaccia
) (Message, error) {
    var orig Message
    var contentStr string

    // 1) Recupera il messaggio originale
    err := db.c.QueryRow(
        `SELECT id, message_content, timestamp, sender_id
           FROM messages
          WHERE id = ?`,
        messageId,
    ).Scan(
        &orig.ID,
        &contentStr,
        &orig.Timestamp,
        &orig.SenderID,                   // string, ok per scan
    )
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

    // 3) (Opzionale) modifica del contenuto per il forward
    forwardedContent := orig.MessageContent

    // 4) Serializzo di nuovo il MessageContent in JSON
    forwardBytes, err := json.Marshal(forwardedContent)
    if err != nil {
        return orig, err
    }

    // 5) Inserimento nella conversazione di destinazione,
    //    convertendo senderID in stringa
    now := time.Now()
    res, err := db.c.Exec(
        `INSERT INTO messages (conversation_id, message_content, timestamp, sender_id)
         VALUES (?, ?, ?, ?)`,
        targetConversationId,
        string(forwardBytes),
        now,
        strconv.FormatUint(senderID, 10), // conv uint64->string
    )
    if err != nil {
        return orig, err
    }

    newID, err := res.LastInsertId()
    if err != nil {
        return orig, err
    }

    // 6) Costruisco il Message inoltrato
    forwardedMsg := Message{
        ID:             int(newID),
        Timestamp:      now,
        SenderID:       strconv.FormatUint(senderID, 10),
        MessageContent: forwardedContent,
        // Preview, Comments, MessageStatus lasciati vuoti
    }

    return forwardedMsg, nil
}

func (db *appdbimpl) DeleteMessage(conversationID, messageID string, senderID uint64) error {
    // opzionalmente verifica che il senderID corrisponda
    res, err := db.c.Exec(
        `DELETE FROM messages WHERE id = ? AND conversation_id = ? AND sender_id = ?`,
        messageID, conversationID, senderID,
    )
    if err != nil {
        return err
    }
    n, _ := res.RowsAffected()
    if n == 0 {
        return ErrMessageDoesNotExist
    }
    return nil
}

