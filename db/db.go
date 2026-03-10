package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", "./bunkerchat.db")
	if err != nil {
		log.Fatal("Failed to connect to Database:", err)
		return
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database Unreachable:", err)
	}
	log.Println("Database Connected")
	createTables()
}

func createTables() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	
	);`

	messageTable := `
	CREATE TABLE IF NOT EXISTS messages(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sent_from_id INTEGER NOT NULL,
		sent_to_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		read BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sent_from_id) REFERENCES users(id),
		FOREIGN KEY (sent_to_id) REFERENCES users(id)
	);`

	_, err := DB.Exec(userTable)
	if err != nil {
		log.Fatal("Failed to create users Table:", err)
	}

	_, err = DB.Exec(messageTable)
	if err != nil {
		log.Fatal("Failed to create message Table:", err)
	}
	log.Println("Tables ready")
}
