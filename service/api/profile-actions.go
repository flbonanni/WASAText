package api

import (
	"encoding/json"
	"net/http"

	"github.com/flbonanni/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"github.com/flbonanni/WASAText/service/database"
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
    // 1. Estrai l’user ID dal token
    tokenID := getToken(r.Header.Get("Authorization"))

    // 2. Verifica che l’utente esista davvero
    dbUser, err := rt.db.CheckUserById(database.User{ID: tokenID})
    if err != nil {
        http.Error(w, "User does not exist", http.StatusUnauthorized)
        return
    }

    // 3. Estrai la foto dal multipart form
    file, header, err := r.FormFile("photo")
    if err != nil {
        http.Error(w, "Invalid photo upload: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    // 4. Leggi tutto in un buffer
    buf := &bytes.Buffer{}
    if _, err := io.Copy(buf, file); err != nil {
        http.Error(w, "Failed to read photo: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // 5. Costruisci l’oggetto Photo (senza campo Filename, e con Date)
	photo := database.Photo{
       UserId: dbUser.ID,
       File:   buf.Bytes(),
       Date:   time.Now().Format(time.RFC3339),
   	}

    // 6. Salva la foto
    if err := rt.db.ChangeUserPhoto(dbUser, photo); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(photo)
}