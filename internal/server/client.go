package server

import (
	"context"
	"encoding/json"
	"log"

	"github.com/coder/websocket"
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	hub        *Hub
	room       string // The room the client is currently in
	id         string // Unique player identifier
	// egress is used to avoid concurrent writes to the websocket connection.
	egress chan Event
}

func (c *Client) writeMessages() {
	defer c.hub.removeClient(c)

	for message := range c.egress {
		data, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			continue
		}

		if err := c.connection.Write(context.Background(), websocket.MessageText, data); err != nil {
			log.Printf("Failed to write message: %v", err)
		}
	}
	// Channel closed, notify client
	// if err := c.connection.Write(context.Background(), websocket.MessageText, []byte("No eager, connection closed")); err != nil {
	// 	log.Printf("Failed to write message: %v", err)
	// }
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		connection: conn,
		hub:        hub,
		egress:     make(chan Event),
	}
}

// generateRandomID returns a random string to be used as a client ID.
// func generateRandomID(length int) string {
// 	charset := "abcdefghijklmnopqrstuvwxyz"
// 	b := make([]byte, length)
// 	for i := range b {
// 		b[i] = charset[rand.Intn(len(charset))]
// 	}
// 	return string(b)
// }

func (c *Client) readMessages() {
	defer c.hub.removeClient(c)
	c.connection.SetReadLimit(1024) // Set a read limit to prevent large messages
	for {
		_, payload, err := c.connection.Read(context.Background())
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure || websocket.CloseStatus(err) == websocket.StatusGoingAway {
				log.Printf("Connection closed normally")
			} else {
				log.Printf("Failed to read message: %v", err)
			}
			break
		}

		var event Event
		if err := json.Unmarshal(payload, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			continue
		}

		if err := c.hub.routeEvent(event, c); err != nil {
			log.Printf("Failed to route event: %v", err)
			continue
		}
	}
}
