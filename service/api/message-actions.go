package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/flbonanni/WASAText/service/api/reqcontext"
	"github.com/flbonanni/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)


func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var conversation Conversation
	var message database.Message
	// estrarre un token dall'header
	token := getToken(r.Header.Get("Authorization"))
	user.ID = token
	dbUser, err := rt.db.CheckUserById(user.ToDatabase())
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
	}
	user.FromDatabase(dbUser)

	// Get the user's conversations from the database
	conversation, err = rt.db.GetConversation(conversation.ConversationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	conversation.ConvFromDatabase(conversation)

	// Decodifica il messaggio inviato nel body della richiesta
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Imposta il timestamp corrente
	message.Timestamp = time.Now()

	// Salva il messaggio nel database associandolo alla conversazione
	savedMessage, err := rt.db.SendMessage(conversation.ConversationID, message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Aggiorna il messaggio con i dati restituiti dal DB (ad esempio, un ID generato)
	message = savedMessage

	// Respond with the conversations
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(message)
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Verifica autenticazione
	var user User
	token := getToken(r.Header.Get("Authorization"))
    requestUser.ID = token
    user, err := rt.db.CheckUserById(requestUser.ToDatabase())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbUser)

	// Estrai parametri dalla URL
	_ = ps.ByName("username")              // username del richiedente (per eventuali controlli)
	// conversationId := ps.ByName("conversation_id")
	messageId := ps.ByName("message_id")

	// Decodifica il body in una mappa senza definire una nuova struct
	var reqBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Recupera i valori richiesti dalla mappa
	targetConversationId, ok := reqBody["target_conversation_id"]
	if !ok || targetConversationId == "" {
		http.Error(w, "target_conversation_id mancante", http.StatusBadRequest)
		return
	}
	// recipient_username è opzionale se non fornito
	recipientUsername := reqBody["recipient_username"]

	// Esegui il forward del messaggio tramite il layer DB (funzione ipotetica)
	forwardedMsg, err := rt.db.ForwardMessage(messageId, targetConversationId, recipientUsername, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Rispondi con il messaggio inoltrato (HTTP 200)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(forwardedMsg)
}


func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Verifica autenticazione
	var user User
	token := getToken(r.Header.Get("Authorization"))
	dbUser, err := rt.db.CheckUserById(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbUser)

	// Estrai parametri dalla URL
	_ = ps.ByName("username")            // username del richiedente (per eventuali controlli)
	conversationId := ps.ByName("conversation_id")
	messageId := ps.ByName("message_id")

	// Chiamata al layer DB per eliminare il messaggio (funzione ipotetica)
	err = rt.db.DeleteMessage(conversationId, messageId, user.ID)
	if err != nil {
		// Gestisci eventuali errori (es. 404 per messaggio non trovato, 403 per permessi insufficienti)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Risposta senza contenuto (HTTP 204)
	w.WriteHeader(http.StatusNoContent)
}