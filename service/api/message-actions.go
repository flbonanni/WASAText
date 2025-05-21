package api

import (
    "encoding/json"
    "errors"
    "net/http"
    "strconv"
    "time"

	"github.com/flbonanni/WASAText/service/api/reqcontext"
	"github.com/flbonanni/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
)


func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    // 1) Autenticazione
    token := getToken(r.Header.Get("Authorization"))
    user := User{ID: token}
    dbUser, err := rt.db.CheckUserById(user.ToDatabase())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    user.FromDatabase(dbUser)

    // 2) Decodifica del body JSON
    var payload struct {
        Type         string   `json:"type"`
        Content      string   `json:"content"`
        Participants []string `json:"participants,omitempty"`
    }
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 3) Recupero o creazione conversazione
    conversationID := ps.ByName("conversation_id")
    conv, err := rt.db.GetConversation(conversationID)
    if err != nil {
        if errors.Is(err, database.ErrConversationDoesNotExist) {
            if len(payload.Participants) < 2 {
                http.Error(w,
                    "conversation does not exist; provide at least two participants to create it",
                    http.StatusBadRequest)
                return
            }
            conv, err = rt.db.CreateConversation(conversationID, payload.Participants)
            if err != nil {
                http.Error(w, "cannot create conversation: "+err.Error(),
                    http.StatusInternalServerError)
                return
            }
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    // 4) Costruzione del messaggio
    var msg database.Message
    msg.Timestamp = time.Now()
    // Converto user.ID (uint64) a string per SenderID
    msg.SenderID = strconv.FormatUint(user.ID, 10)

    switch payload.Type {
    case "text":
        msg.MessageContent = database.MessageContent{
            Type: payload.Type,
            Text: payload.Content,
        }
    case "image":
        msg.MessageContent = database.MessageContent{
            Type:     payload.Type,
            ImageURL: payload.Content,
        }
    default:
        http.Error(w, "unsupported message type", http.StatusBadRequest)
        return
    }

    // 5) Salvataggio nel DB
    msgSaved, err := rt.db.SendMessage(conv.ConversationID, msg)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 6) Risposta JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(msgSaved); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Verifica autenticazione
	var user User
	var requestUser User
	token := getToken(r.Header.Get("Authorization"))
    requestUser.ID = token
    dbUser, err := rt.db.CheckUserById(requestUser.ToDatabase())

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
    // 1) Autenticazione
    token := getToken(r.Header.Get("Authorization"))
    user := User{ID: token}
    dbUser, err := rt.db.CheckUserById(user.ToDatabase())
    if err != nil {
        http.Error(w, "User does not exist", http.StatusUnauthorized)
        return
    }
    user.FromDatabase(dbUser)

    // 2) Parametri URL
    // username := ps.ByName("username")
    // (opzionale) potresti voler controllare che username == user.Username

    conversationID := ps.ByName("conversation_id")
    messageID := ps.ByName("message_id")

    // 3) Chiamata al DB
    if err := rt.db.DeleteMessage(conversationID, messageID, token); err != nil {
        if err == database.ErrMessageDoesNotExist {
            http.Error(w, "Message not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // 4) Risposta
    w.WriteHeader(http.StatusNoContent) // 204
}
