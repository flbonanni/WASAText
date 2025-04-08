package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/flbonanni/WASAText/service/api/reqcontext"
	"github.com/flbonanni/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Verifica autenticazione
	var user User
	token := getToken(r.Header.Get("Authorization"))
	user.ID = token
	dbUser, err := rt.db.CheckUserById(user.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbUser)

	// Estrai parametri dalla URL
	_ = ps.ByName("username")            // eventuale controllo sulla coerenza con l'utente loggato
	conversationId := ps.ByName("conversation_id")
	messageId := ps.ByName("message_id")

	// Decodifica il body in una mappa per ottenere l'emoji
	var reqBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	emoji, ok := reqBody["emoji"]
	if !ok || emoji == "" {
		http.Error(w, "emoji mancante", http.StatusBadRequest)
		return
	}

	// Aggiungi l'emoji reaction al messaggio nel database
	// (La funzione rt.db.CommentMessage è ipotetica e deve gestire l'associazione tra
	//  conversationId, messageId, emoji e l'ID dell'utente che aggiunge il commento)
	err = rt.db.CommentMessage(conversationId, messageId, emoji, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Rispondi con un messaggio di conferma
	response := map[string]string{
		"message": "Emoji reaction added successfully.",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Verifica autenticazione
	var user User
	// token := getToken(r.Header.Get("Authorization"))
	username := ps.ByName("username")
	dbuser, err := rt.db.GetUserId(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)

	// Estrai parametri dalla URL
	_ = ps.ByName("username")            // eventuale controllo sull'utente loggato
	conversationId := ps.ByName("conversation_id")
	messageId := ps.ByName("message_id")

	// Rimuove l'emoji reaction dal messaggio nel database
	err = rt.db.UncommentMessage(conversationId, messageId, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Rispondi con HTTP 204 No Content
	w.WriteHeader(http.StatusNoContent)
}