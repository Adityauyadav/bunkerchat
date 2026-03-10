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
	query := `SELECT id, username, password_hash, token, created_at FROM users WHERE username = ?`
	row := DB.QueryRow(query, username)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Token, &user.CreatedAt)
	return user, err
}

func SaveToken(userID int, token string) error {
	query := `UPDATE users SET token = ? WHERE id = ?`
	_, err := DB.Exec(query, token, userID)
	return err
}

func ClearToken(userID int) error {
	query := `UPDATE users SET token = NULL WHERE id = ? `
	_, err := DB.Exec(query, userID)
	return err

}
