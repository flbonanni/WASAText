package database

import (
	"database/sql"
	"github.com/flbonanni/WASAText/datamodels"
)

func (db *appdbimpl) CreateUser(u datamodels.User) (datamodels.User, error) {
	res, err := db.c.Exec("INSERT INTO users(username) VALUES (?)", u.CurrentUsername)
	if err != nil {
		// Se l'INSERT fallisce, prova a recuperare l'utente esistente
		var user datamodels.User
		err = db.c.QueryRow(`SELECT id, username FROM users WHERE username = ?`, u.CurrentUsername).Scan(&user.ID, &user.CurrentUsername)
		if err != nil {
			if err == sql.ErrNoRows {
				// Se l'utente non esiste, restituisci l'errore originale dell'INSERT
				return datamodels.User{}, err
			}
			// Altri errori di QueryRow devono essere gestiti
			return datamodels.User{}, err
		}
		// Se l'utente esiste già, restituiscilo senza errore
		return user, nil
	}

	// Recupera l'ID dell'utente appena inserito
	lastInsertID, err := res.LastInsertId()
	if err != nil {
    return datamodels.User{}, err
	}

	// CONVERSIONE DA int64 A int
	u.ID = int(lastInsertID) 
	return u, nil
}

func (db *appdbimpl) SetUsername(u datamodels.User, username string) (datamodels.User, error) {
	res, err := db.c.Exec(`UPDATE users SET Username=? WHERE Id=? AND Username=?`, u.CurrentUsername, u.ID, username)
	if err != nil {
		return u, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return u, err
	} else if affected == 0 {
		return u, err
	}
	return u, nil
}

func (db *appdbimpl) GetUserId(username string) (datamodels.User, error) {
	var user datamodels.User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE username = ?`, username).Scan(&user.ID, &user.CurrentUsername); err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
	}
	return user, nil
}

func (db *appdbimpl) CheckUserByUsername(u datamodels.User) (datamodels.User, error) {
	var user datamodels.User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE username = ?`, u.CurrentUsername).Scan(&user.ID, &user.CurrentUsername); err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
	}
	return user, nil
}

func (db *appdbimpl) CheckUserById(u User) (User, error) {
	var user datamodels.User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE id = ?`, u.Id).Scan(&user.Id, &user.Username); err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUserDoesNotExist
		}
	}
	return user, nil
}

func (db *appdbimpl) CheckUser(u datamodels.User) (datamodels.User, error) {
	var user datamodels.User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE id = ? AND username = ?`, u.ID, u.CurrentUsername).Scan(&user.ID, &user.CurrentUsername); err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
	}
	return user, nil
}