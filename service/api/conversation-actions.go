package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)


func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	// estrarre un token dall'header
	token := getToken(r.Header.Get("Authorization"))
	dbUser, err := rt.db.CheckUserById(token)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
	}
	user.FromDatabase(dbUser)

	username := ps.ByName("username")

	// Get the user's conversations from the database
	conversations, err := rt.db.GetConversations(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the conversations
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(conversations)
}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var conversation Conversation
	// estrarre un token dall'header
	token := getToken(r.Header.Get("Authorization"))
	dbUser, err := rt.db.CheckUserById(token)
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
	}
	user.FromDatabase(dbUser)


	conversationId := ps.ByName("conversation_id")

	// Get the user's conversation from the database
	conversation, err := rt.db.GetConversation(ConversationId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the conversation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(conversation)
}