package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coder/websocket"
)

type Hub struct {
	client ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func (h *Hub) setupEventHandlers() {
	h.handlers[EventSendMessage] = SendMessage
	h.handlers[EventJoinRoom] = JoinRoomHandler
	h.handlers[EventCreateRoom] = SendMessage
	h.handlers[EventLeaveRoom] = SendMessage
	h.handlers[EventListRooms] = SendMessage
}

func JoinRoomHandler(event Event, c *Client) error {
	var joinRoomEvent JoinRoomEvent
	if err := json.Unmarshal(event.Payload, &joinRoomEvent); err != nil {
		return fmt.Errorf("failed to unmarshal join room event: %v", err)
	}
	c.room = joinRoomEvent.Room
	return nil
}

func SendMessage(event Event, c *Client) error {
	var chatEvent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		log.Printf("Failed to unmarshal send message event: %v", err)
		return err
	}

	var broadMessage NewMessageEvent
	broadMessage.Sent = time.Now()
	broadMessage.From = chatEvent.From
	broadMessage.Message = chatEvent.Message

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("Failed to marshal new message event: %v", err)
	}

	outgoingEvent := Event{
		Type:    EventNewMessage,
		Payload: data,
	}
	for client := range c.hub.client {
		if client.room == c.room {
			client.egress <- outgoingEvent
		}
	}

	return nil
}

func NewHub() *Hub {
	h := &Hub{
		client:   make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	h.setupEventHandlers()
	return h
}

func (h *Hub) routeEvent(event Event, c *Client) error {
	if handler, ok := h.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("no handler for event type: %s", event.Type)
}

func (h *Hub) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"}, // Allow all origins for development
	})
	if err != nil {
		log.Printf("Failed to accept websocket connection: %v", err)
		return
	}
	//defer conn.Close(websocket.StatusInternalError, "Connection closed")

	client := NewClient(conn, h)
	h.addClient(client)

	// Handle websocket messages in this goroutine to keep connection open
	go client.readMessages()
	go client.writeMessages()
}

func (h *Hub) addClient(client *Client) {
	h.Lock()
	defer h.Unlock()
	h.client[client] = true
}

func (h *Hub) removeClient(client *Client) {
	h.Lock()
	defer h.Unlock()
	if _, ok := h.client[client]; ok {
		client.connection.Close(websocket.StatusNormalClosure, "Connection closed normally")
		delete(h.client, client)
	}
}
