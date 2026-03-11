package main

import (
	"encoding/json"
	"log"

	"github.com/adityauyadav/bunkerchat/db"
	"github.com/adityauyadav/bunkerchat/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	username    string
	userID      int
	recipientID int
	recipient   string
	send        chan []byte
}

func (c *Client) outgoing() {
	defer c.hub.Unregister(c.username)
	defer c.conn.Close()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var packets models.MessagePacket

		err = json.Unmarshal(message, &packets)
		if err != nil {
			log.Println("Invalid Message Format:", err)
			continue
		}
		db.SaveMessage(c.userID, c.recipientID, packets.Content)
		c.hub.SendMessage(c.recipient, message)
	}

}

func (c *Client) incoming() {
	for {
		message, ok := <-c.send
		if !ok {
			break
		}
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
}
