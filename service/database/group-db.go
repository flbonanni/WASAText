package database

import (
	"fmt"
	"database/sql"
	"strings"
	"mime/multipart"
	"io"
	"log"
	"time"
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

func (db *appdbimpl) UpdateGroupPhoto(groupId string, adminID uint64, photoData multipart.File) error {
    defer photoData.Close()

    // 1) Leggi i byte del file
    photoBytes, err := io.ReadAll(photoData)
    if err != nil {
        return err
    }

    // 2) Esegui l'UPDATE SOLO su group_id
    res, err := db.c.Exec(
        `UPDATE groups
            SET photo = ?
          WHERE group_id = ?`,
        photoBytes,
        groupId,
    )
    if err != nil {
        return err
    }

    // 3) Controlla quante righe sono state modificate
    affected, err := res.RowsAffected()
	log.Printf("Trying to update group %s with admin ID %d", groupId, adminID)
	log.Printf("Rows affected: %d", affected)
    if affected == 0 {
        return ErrGroupNotUpdated
    }

    return nil
}

func (db *appdbimpl) CreateGroup(
    adminID uint64,
    groupName string,
    description string,
    members []string,
) (string, error) {
    // 1) Genera un ID univoco per il gruppo
    //    Qui usiamo un prefisso + timestamp UNIX, ma puoi sostituire con uuid.New().String()
    groupID := fmt.Sprintf("group%d", time.Now().UnixNano())

    // 2) Prepara la stringa dei membri
    membersStr := strings.Join(members, ",")

    // 3) Esegui l'INSERT con tutti e 5 i placeholder
    _, err := db.c.Exec(
        `INSERT INTO groups (group_id, admin_id, group_name, description, members)
         VALUES (?, ?, ?, ?, ?)`,
        groupID,
        adminID,
        groupName,
        description,
        membersStr,
    )
    if err != nil {
        return "", err
    }

    // 4) Ritorna il nuovo groupID
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