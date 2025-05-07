package database

import (
	"database/sql"
)

func (db *appdbimpl) CreateUser(u User) (User, error) {
	res, err := db.c.Exec("INSERT INTO users(username) VALUES (?)", u.CurrentUsername)
	if err != nil {
		// Se l'INSERT fallisce, prova a recuperare l'utente esistente
		var user User
		err = db.c.QueryRow(`SELECT id, username FROM users WHERE username = ?`, u.CurrentUsername).Scan(&user.ID, &user.CurrentUsername)
		if err != nil {
			if err == sql.ErrNoRows {
				// Se l'utente non esiste, restituisci l'errore originale dell'INSERT
				return User{}, err
			}
			// Altri errori di QueryRow devono essere gestiti
			return User{}, err
		}
		// Se l'utente esiste già, restituiscilo senza errore
		return user, nil
	}

	// Recupera l'ID dell'utente appena inserito
	lastInsertID, err := res.LastInsertId()
	if err != nil {
    return User{}, err
	}

	// CONVERSIONE DA int64 A int
	u.ID = uint64(lastInsertID)
	return u, nil
}

/* func (db *appdbimpl) SetUsername(u User, username string) (User, error) {
	res, err := db.c.Exec(`UPDATE users SET Username=? WHERE Id=? AND Username=?`, username, u.ID, u.CurrentUsername)
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
} */

func (db *appdbimpl) SetUsername(u User, oldUsername string) (User, error) {
	// 1) esegui l'UPDATE con i parametri corretti
	res, err := db.c.Exec(
	  `UPDATE users SET Username=? WHERE Id=? AND Username=?`,
	  u.Username,
	  u.ID,
	  oldUsername,
	)
	if err != nil {
	  return u, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
	  return u, fmt.Errorf("nessuna riga aggiornata: id=%d username=%q", u.ID, oldUsername)
	}
	// 2) recupera il record aggiornato
	row := db.c.QueryRow(`SELECT Id, Username FROM users WHERE Id = ?`, u.ID)
	var updated User
	if err := row.Scan(&updated.ID, &updated.Username); err != nil {
	  return u, err
	}
	return updated, nil
  }
  

func (db *appdbimpl) GetUserId(username string) (User, error) {
	var user User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE username = ?`, username).Scan(&user.ID, &user.CurrentUsername); err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
	}
	return user, nil
}

func (db *appdbimpl) CheckUserByUsername(u User) (User, error) {
	var user User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE username = ?`, u.CurrentUsername).Scan(&user.ID, &user.CurrentUsername); err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
	}
	return user, nil
}

func (db *appdbimpl) CheckUserById(u User) (User, error) {
	var user User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE id = ?`, u.ID).Scan(&user.ID, &user.CurrentUsername); err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUserDoesNotExist
		}
	}
	return user, nil
}

func (db *appdbimpl) CheckUser(u User) (User, error) {
	var user User
	if err := db.c.QueryRow(`SELECT id, username FROM users WHERE id = ? AND username = ?`, u.ID, u.CurrentUsername).Scan(&user.ID, &user.CurrentUsername); err != nil {
		if err == sql.ErrNoRows {
			return user, err
		}
	}
	return user, nil
}