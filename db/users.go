package db

import (
	"github.com/adityauyadav/bunkerchat/models"
)

func CreateUser(username string, passwordhash string) error {
	query := `INSERT INTO users (username, password_hash) VALUES (?,?)`
	_, err := DB.Exec(query, username, passwordhash)
	return err

}

func GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password_hash, created_at FROM users WHERE username = ?`
	row := DB.QueryRow(query, username)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	return user, err
}

func GetUserByID(id int) (*models.User, error) {
	query := `SELECT id, username, password_hash, created_at FROM users WHERE id = ?`
	row := DB.QueryRow(query, id)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt)
	return user, err
}
