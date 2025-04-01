package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/flbonanni/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	// errore = decodifica del corpo di r + popolamento della var user
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
	var user User
	var requestUser User
	var profile Profile
	// estrarre un token dall'header
	token := getToken(r.Header.Get("Authorization"))
	requestUser.Id = token
	// controlla che esista un tale utente
	dbrequestuser, err := rt.db.CheckUserById(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// popola l'utente richiedente
	requestUser.FromDatabase(dbrequestuser)
	// popola lo username dal parametro
	username := ps.ByName("username")
	
	// stessa cosa di prima per identificare l'utente nel db
	dbuser, err := rt.db.GetUserId(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)

	// aggiornamento di profile con le variabili richieste
	profile.RequestId = token
	profile.Id = user.Id
	profile.Username = user.Username
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(profile)
}

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
	user.Id = token

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