package main

import (
	"encoding/json"
	"net/http"

	"github.com/adityauyadav/bunkerchat/auth"
	"github.com/adityauyadav/bunkerchat/db"
	"github.com/gorilla/websocket"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	req := RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return

	}
	err = db.CreateUser(req.Username, hash)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user created successfully",
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	req := LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	user, err := db.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	validation := auth.CheckPassword(req.Password, user.PasswordHash)
	if !validation {
		http.Error(w, "Invalid Password", http.StatusBadRequest)
		return
	}
	token, err := auth.GenerateToken(user.ID, req.Username)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})

}

func wsHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.URL.Query().Get("token")
		recipient := r.URL.Query().Get("recipient")

		if token == "" {
			http.Error(w, "Missing token", http.StatusBadRequest)
			return
		}

		userID, username, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusBadRequest)
			return
		}

		recipientUser, err := db.GetUserByUsername(recipient)
		if err != nil {
			http.Error(w, "User Doesn't Exist", http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade connection", http.StatusBadRequest)
			return
		}
		client := &Client{
			hub:         hub,
			conn:        conn,
			send:        make(chan []byte, 256),
			userID:      userID,
			username:    username,
			recipient:   recipientUser.Username,
			recipientID: recipientUser.ID,
		}

		client.hub.Register(client.username, client.conn)

		go client.incoming()
		go client.outgoing()
	}
}
