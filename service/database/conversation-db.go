package database

import (
	"database/sql"
	"errors"
	"strings"
)

var ErrConversationDoesNotExist = errors.New("conversation does not exist")

func (db *appdbimpl) GetConversations(username string) ([]api.Conversation, error) {
	rows, err := db.c.Query(
		`SELECT conversation_id, participants, last_message FROM conversations 
		 WHERE FIND_IN_SET(?, participants) > 0`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []api.Conversation
	for rows.Next() {
		var conv api.Conversation
		var participantsStr string
		if err := rows.Scan(&conv.ConversationID, &participantsStr, &conv.LastMessage); err != nil {
			return nil, err
		}
		// Converti la stringa dei partecipanti in slice (assumendo separazione tramite virgola)
		conv.Participants = strings.Split(participantsStr, ",")
		conversations = append(conversations, conv)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return conversations, nil
}

func (db *appdbimpl) GetConversation(conversationId string) (api.Conversation, error) {
	var conv api.Conversation
	var participantsStr string
	if err := db.c.QueryRow(
		`SELECT conversation_id, participants, last_message FROM conversations 
		 WHERE conversation_id = ?`, conversationId).Scan(&conv.ConversationID, &participantsStr, &conv.LastMessage); err != nil {
		if err == sql.ErrNoRows {
			return conv, ErrConversationDoesNotExist
		}
		return conv, err
	}
	conv.Participants = strings.Split(participantsStr, ",")
	return conv, nil
}
