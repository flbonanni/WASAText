package database

import (
	"database/sql"
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

func (db *appdbimpl) ChangeUserPhoto(u User, photo Photo) error {
    // Salva il BLOB file e il nome file nel campo `photo_filename`
    _, err := db.c.Exec(
        `UPDATE users SET photo = ?, photo_filename = ? WHERE id = ?`,
        photo.File, photo.Filename, u.ID,
    )
    if err != nil {
        return err
    }
    return nil
}
