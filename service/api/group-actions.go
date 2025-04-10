package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/flbonanni/WASAText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    token := getToken(r.Header.Get("Authorization"))
    var requestUser User
    requestUser.ID = token
    user, err := rt.db.CheckUserById(requestUser.ToDatabase())
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    groupId := ps.ByName("group_id")
	// ogni campo è una coppia chiave/valore di tipo stringa
    var reqBody map[string]string
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

	// Estrae il valore associato alla chiave "name" dalla mappa.
	// Verifica che il campo esista (ok deve essere true) e che la lunghezza del nome sia compresa tra 3 e 50 caratteri.
    groupName, ok := reqBody["name"]
    if !ok || len(groupName) < 3 || len(groupName) > 50 {
        http.Error(w, "Invalid group name", http.StatusBadRequest)
        return
    }

    if err := rt.db.UpdateGroupName(groupId, user.ID, groupName); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Group name updated successfully."})
}

func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var requestuser User
	token := getToken(r.Header.Get("Authorization"))
	requestuser.ID = token

	user, err := rt.db.CheckUserById(requestuser.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = r.ParseMultipartForm(10 << 20) 
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	groupId := ps.ByName("group_id")
	file, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = rt.db.UpdateGroupPhoto(groupId, user.ID, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Group picture uploaded successfully."})
}

func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    var requestuser User
	token := getToken(r.Header.Get("Authorization"))
	requestuser.ID = token
	user, err := rt.db.CheckUserById(requestuser.ToDatabase())	
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }
    
    var reqBody struct {
        GroupName   string   `json:"group_name"`
        Description string   `json:"description"`
        Members     []string `json:"members"`
    }
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    groupId, err := rt.db.CreateGroup(user.ID, reqBody.GroupName, reqBody.Description, append(reqBody.Members, user.ID))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"group_id": groupId, "group_name": reqBody.GroupName})
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    var user User
	token := getToken(r.Header.Get("Authorization"))
	user.ID = token
	user, err := rt.db.CheckUserById(user.ToDatabase())
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }
    
    groupId := ps.ByName("group_id")
    var reqBody map[string]string
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    newMember, ok := reqBody["new_member_username"]
    if !ok || len(newMember) < 3 || len(newMember) > 30 {
        http.Error(w, "Invalid member username", http.StatusBadRequest)
        return
    }
    
    err = rt.db.AddMemberToGroup(groupId, user.ID, newMember)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Member added successfully."})
}

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    var user User
	token := getToken(r.Header.Get("Authorization"))
	user.ID = token
	user, err := rt.db.CheckUserById(user.ToDatabase())
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }
    
    groupId := ps.ByName("group_id")
    memberUsername := ps.ByName("member_username")
    if memberUsername != user.CurrentUsername {
        http.Error(w, "Unauthorized action", http.StatusForbidden)
        return
    }
    
    err = rt.db.RemoveMemberFromGroup(groupId, memberUsername)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusNoContent)
}
