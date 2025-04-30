package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/flbonanni/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)


func (rt *_router) getUserPicture(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract the username from the URL
	username := ps.ByName("username")

	// Get the user's profile picture from the database
	picture, err := rt.db.GetUserPicture(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the image
	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(picture)
}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user database.User
	var photo database.Photo
	token := getToken(r.Header.Get("Authorization"))
	user.ID = token
	dbUser, err := rt.db.CheckUserById(user.ToDatabase())
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
	}
	user.FromDatabase(dbUser)

	err = json.NewDecoder(r.Body).Decode(&photo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	// change the user photo
	err = rt.db.ChangeUserPhoto(user, photo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(photo)
}