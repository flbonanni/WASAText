package api

import (
	"encoding/json"
	"net/http"

	"github.com/flbonanni/WASAText/service/api/reqcontext"
	"github.com/flbonanni/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)


func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    // 1) Estrai il token ed esegui l’autenticazione
    token := getToken(r.Header.Get("Authorization"))
    user := User{ID: token}
    dbUser, err := rt.db.CheckUserById(user.ToDatabase())
    if err != nil {
        http.Error(w, "User does not exist", http.StatusUnauthorized)
        return
    }
    user.FromDatabase(dbUser)

    // 2) Prendi l’ID della conversazione dai path params
    conversationID := ps.ByName("conversation_id")

    // 3) Recupera la conversazione
    conv, err := rt.db.GetConversation(conversationID)
    if err != nil {
        if err == database.ErrConversationDoesNotExist {
            http.Error(w, "Conversation does not exist", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // 4) Rispondi con la conversazione
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(conv)
}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var conversation database.Conversation
	// estrarre un token dall'header
	//token := getToken(r.Header.Get("Authorization"))
	dbUser, err := rt.db.CheckUserById(user.ToDatabase())
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
	}
	user.FromDatabase(dbUser)


	conversationId := ps.ByName("conversation_id")

	// Get the user's conversation from the database
	conversation, err = rt.db.GetConversation(conversationId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the conversation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(conversation)
}