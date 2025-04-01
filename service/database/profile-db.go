package database

import (
	"database/sql"
	"github.com/flbonanni/WASAText/datamodels"
)

func (db *appdbimpl) GetUserPicture(username string) ([]byte, error) {
	var picture []byte
	// Esegue la query per ottenere la foto (campo photo) dell'utente
	if err := db.c.QueryRow(`SELECT photo FROM users WHERE username = ?`, username).Scan(&picture); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserDoesNotExist
		}
		return nil, err
	}
	return picture, nil
}

func (db *appdbimpl) ChangeUserPhoto(u datamodels.User, photo datamodels.Photo) error {
	// Esegue l'update della foto dell'utente identificato da u.Id
	res, err := db.c.Exec(`UPDATE users SET photo = ? WHERE id = ?`, photo.File, u.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	} else if affected == 0 {
		return ErrUserDoesNotExist
	}
	return nil
}
