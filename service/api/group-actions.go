package api

import (
	"encoding/json"
	"net/http"

	"github.com/flbonanni/WASAText/service/api/reqcontext"
    "github.com/flbonanni/WASAText/service/database"
	"github.com/julienschmidt/httprouter"
    "strconv"
    "fmt"
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
    
    groupId, err := rt.db.CreateGroup(
        user.ID,
        reqBody.GroupName,
        reqBody.Description,
        append(reqBody.Members, strconv.FormatUint(user.ID, 10)),
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"group_id": groupId, "group_name": reqBody.GroupName})
}

func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
    // 1) Estraggo token e costruisco user
    token := getToken(r.Header.Get("Authorization"))
    user := User{ID: token}

    // 2) Verifico che esista nel DB
    dbUser, err := rt.db.CheckUserById(user.ToDatabase())
    if err != nil {
        http.Error(w, "User does not exist", http.StatusUnauthorized)
        return
    }
    user.FromDatabase(dbUser)

    // (Facoltativo) Controlla che il path param username corrisponda
    usernameParam := ps.ByName("username")
    if usernameParam != user.CurrentUsername {
        http.Error(w, "username mismatch", http.StatusForbidden)
        return
    }

    // 3) Leggo il group_id
    groupID := ps.ByName("group_id")

    // 4) Decodifico il body
    var reqBody struct {
        NewMemberUsername string `json:"new_member_username"`
    }
    if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    newMember := reqBody.NewMemberUsername
    if len(newMember) < 3 || len(newMember) > 30 {
        http.Error(w, "Invalid member username", http.StatusBadRequest)
        return
    }

    // 5) Invoco il DB
    if err := rt.db.AddMemberToGroup(groupID, user.ID, newMember); err != nil {
        switch err {
        case database.ErrGroupNotFound:
            http.Error(w, "Group not found", http.StatusNotFound)
        case fmt.Errorf("member already exists"):
            http.Error(w, "Member already exists", http.StatusConflict)
        default:
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // 6) Risposta
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(map[string]string{"message": "Member added successfully."})
}


func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
    groupId := ps.ByName("group_id")
    memberUsername := ps.ByName("member_username")

    // 3) Controlla che sia lo user che se ne vuole andare
    if memberUsername != user.Username {
        http.Error(w, "Unauthorized action", http.StatusForbidden)
        return
    }

    // 4) Rimuovi il membro
    if err := rt.db.RemoveMemberFromGroup(groupId, memberUsername); err != nil {
        switch err {
        case database.ErrGroupNotFound:
            http.Error(w, "Group not found", http.StatusNotFound)
        case fmt.Errorf("member not found in group"):
            http.Error(w, "Member not in group", http.StatusNotFound)
        default:
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // 5) Risposta 204 No Content
    w.WriteHeader(http.StatusNoContent)
}

