package main

import (
	"log"
	"net/http"
	"os"

	"github.com/adityauyadav/bunkerchat/db"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db.Init()

	hub := NewHub()

	registerRoutes(hub)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func registerRoutes(hub *Hub) {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/ws", wsHandler(hub))
	http.HandleFunc("/chat/", chatHistoryHandler)
}
