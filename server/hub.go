package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	mu          sync.Mutex
	connections map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[string]*websocket.Conn),
	}
}

func (h *Hub) Register(username string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connections[username] = conn

}

func (h *Hub) Unregister(username string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.connections, username)

}

func (h *Hub) SendMessage(to string, message []byte) bool {
	h.mu.Lock()
	conn, onlineStatus := h.connections[to]
	h.mu.Unlock()
	if !onlineStatus {
		return false
	}
	conn.WriteMessage(websocket.TextMessage, message)
	return true

}
