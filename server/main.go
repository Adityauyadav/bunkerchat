package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adityauyadav/bunkerchat/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("BunkerChat server starting...")
	db.Init()

	http.HandleFunc("/health", healthHandler)

	log.Println("Server listening on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "BunkerChat server is alive")
}
