package database

import (
	"fmt"
	"database/sql"
	"strings"
)

var (
	ErrGroupNotFound   = fmt.Errorf("group not found")
	ErrGroupNotUpdated = fmt.Errorf("group not updated")
)

// UpdateGroupName aggiorna il nome di un gruppo se l'utente è admin.
func (db *appdbimpl) UpdateGroupName(groupId string, adminID uint64, groupName string) error {
	res, err := db.c.Exec(`UPDATE groups SET group_name = ? WHERE group_id = ? AND admin_id = ?`, groupName, groupId, adminID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrGroupNotUpdated
	}
	return nil
}

// UpdateGroupPhoto aggiorna la foto del gruppo se l'utente è admin.
// Il parametro file è un io.Reader che rappresenta il file caricato.
func (db *appdbimpl) UpdateGroupPhoto(groupId string, adminID uint64, photoData multipart.File) error {
	res, err := db.c.Exec(`UPDATE groups SET photo = ? WHERE group_id = ? AND admin_id = ?`, photoData, groupId, adminID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrGroupNotUpdated
	}
	return nil
}

// CreateGroup crea un nuovo gruppo, impostando l'utente loggato come admin.
// I membri vengono memorizzati come una stringa con valori separati da virgola.
func (db *appdbimpl) CreateGroup(adminID uint64, groupName string, description string, members []string) (string, error) {
	// Converti la slice dei membri in una stringa separata da virgole.
	membersStr := strings.Join(members, ",")
	res, err := db.c.Exec(`INSERT INTO groups (admin_id, group_name, description, members) VALUES (?, ?, ?, ?)`,
		adminID, groupName, description, membersStr)
	if err != nil {
		return "", err
	}
	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	// Genera un groupID (ad esempio "group12345")
	groupID := fmt.Sprintf("group%d", lastInsertID)
	return groupID, nil
}

// AddMemberToGroup aggiunge un nuovo membro a un gruppo esistente.
// Si recupera la lista attuale, si aggiunge il nuovo membro e si aggiorna il record.
func (db *appdbimpl) AddMemberToGroup(groupId string, adminID uint64, newMemberUsername string) error {
	var membersStr string
	err := db.c.QueryRow(`SELECT members FROM groups WHERE group_id = ? AND admin_id = ?`, groupId, adminID).Scan(&membersStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrGroupNotFound
		}
		return err
	}
	members := strings.Split(membersStr, ",")
	// Verifica che il nuovo membro non esista già.
	for _, m := range members {
		if m == newMemberUsername {
			return fmt.Errorf("member already exists")
		}
	}
	members = append(members, newMemberUsername)
	updatedMembersStr := strings.Join(members, ",")
	res, err := db.c.Exec(`UPDATE groups SET members = ? WHERE group_id = ? AND admin_id = ?`, updatedMembersStr, groupId, adminID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrGroupNotUpdated
	}
	return nil
}

// RemoveMemberFromGroup rimuove un membro da un gruppo.
// La funzione aggiorna la lista dei membri rimuovendo il membro specificato.
func (db *appdbimpl) RemoveMemberFromGroup(groupId string, memberUsername string) error {
	var membersStr string
	err := db.c.QueryRow(`SELECT members FROM groups WHERE group_id = ?`, groupId).Scan(&membersStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrGroupNotFound
		}
		return err
	}
	members := strings.Split(membersStr, ",")
	var updatedMembers []string
	found := false
	for _, m := range members {
		if m == memberUsername {
			found = true
			continue
		}
		updatedMembers = append(updatedMembers, m)
	}
	if !found {
		return fmt.Errorf("member not found in group")
	}
	updatedMembersStr := strings.Join(updatedMembers, ",")
	res, err := db.c.Exec(`UPDATE groups SET members = ? WHERE group_id = ?`, updatedMembersStr, groupId)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrGroupNotUpdated
	}
	return nil
}