package api

import (
	"encoding/json"
	"net/http"

	"github.com/flbonanni/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// creazione utente
	dbuser, err := rt.db.CreateUser(user.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// ripopola user con info dal db, ovvero user id + username
	user.FromDatabase(dbuser)

	// risposta in json, risposta positiva, e json del nuovo utente inserito nella risposta w
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrarre il token dall'header Authorization
	token := getToken(r.Header.Get("Authorization"))

	// Popola l'utente richiedente con il token e verifica in DB
	var requestUser User
	requestUser.ID = token
	dbRequestUser, err := rt.db.CheckUserById(requestUser.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	requestUser.FromDatabase(dbRequestUser)

	// Ottieni lo username dal path parameter
	username := ps.ByName("username")

	// Recupera l'utente target dal DB tramite username
	dbUser, err := rt.db.GetUserId(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var targetUser User
	targetUser.FromDatabase(dbUser)

	// Compila il profilo da restituire
	profile := Profile{
		RequestID: requestUser.ID,
		ID:    targetUser.ID,
		Username:  targetUser.CurrentUsername,
		// Aggiungi altri campi di Profile se necessario
	}

	// Invia la risposta JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
	}
}



// func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
// 	var user User
// 	var requestUser User
// 	var profile Profile
// 	// estrarre un token dall'header
// 	token := getToken(r.Header.Get("Authorization"))
// 	requestUser.ID = token
// 	// controlla che esista un tale utente
// 	dbrequestuser, err := rt.db.CheckUserById(requestUser.ToDatabase())
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	// popola l'utente richiedente
// 	requestUser.FromDatabase(dbrequestuser)
// 	// popola lo username dal parametro
// 	username := ps.ByName("username")
	
// 	// stessa cosa di prima per identificare l'utente nel db
// 	dbuser, err := rt.db.GetUserId(username)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	user.FromDatabase(dbuser)

// 	// aggiornamento di profile con le variabili richieste
// 	profile.RequestId = token
// 	profile.Id = user.Id
// 	profile.Username = user.Username
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
	
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	_ = json.NewEncoder(w).Encode(profile)
// }

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	// estrarre username 
	username := ps.ByName("username")
	// controllo se esiste utente
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// estrarre un token dall'header
	token := getToken(r.Header.Get("Authorization"))
	user.ID = token

	// impostare il nuovo username
	dbuser, err := rt.db.SetUsername(user.ToDatabase(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)

	// scrivere risposta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}