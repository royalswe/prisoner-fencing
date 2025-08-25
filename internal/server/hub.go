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
	h.handlers[EventLeaveRoom] = SendMessage
	h.handlers[EventListRooms] = ListRoomHandler
	h.handlers[EventInitClient] = InitClientHandler
	h.handlers[EventGameAction] = GameActionHandler
}

func InitClientHandler(event Event, c *Client) error {
	var initEvent InitClientEvent
	if err := json.Unmarshal(event.Payload, &initEvent); err != nil {
		return fmt.Errorf("failed to unmarshal init client event: %v", err)
	}
	c.id = initEvent.PlayerId

	// Notify the client of their player ID
	emit(Event{
		Type:    EventInitClient,
		Payload: json.RawMessage(fmt.Sprintf(`{"playerId": "%s"}`, c.id)),
	}, c)

	return nil
}

func ListRoomHandler(event Event, c *Client) error {
	// var listRoomEvent ListRoomEvent
	// if err := json.Unmarshal(event.Payload, &listRoomEvent); err != nil {
	// 	return fmt.Errorf("failed to unmarshal list room event: %v", err)
	// }

	var rooms []string
	for client := range c.hub.client {
		if client.room != "" && !contains(rooms, client.room) {
			rooms = append(rooms, client.room)
		}
	}
	fmt.Printf("Available rooms: %v\n", rooms)
	response := ListRoomResponse{Rooms: rooms}
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal list room response: %v", err)
	}

	outgoingEvent := Event{
		Type:    EventListRooms,
		Payload: data,
	}
	roomEmit(outgoingEvent, "", c.hub)

	return nil
}

// contains checks if a slice of strings contains a specific string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func JoinRoomHandler(event Event, c *Client) error {
	var joinRoomEvent JoinRoomEvent
	if err := json.Unmarshal(event.Payload, &joinRoomEvent); err != nil {
		return fmt.Errorf("failed to unmarshal join room event: %v", err)
	}
	c.room = joinRoomEvent.Room

	// Initialize GameState for the room if it doesn't exist
	if _, ok := roomStates[c.room]; !ok {
		roomStates[c.room] = &GameState{
			Turn:         1,
			MaxTurns:     20,
			PlayerStates: make(map[string]PlayerState),
		}
	}
	gs := roomStates[c.room]
	// If less than two players, add this client as PlayerState
	if _, exists := gs.PlayerStates[c.id]; !exists && len(gs.PlayerStates) < 2 {
		pos := 2
		if len(gs.PlayerStates) > 0 {
			pos = 4
		}
		gs.PlayerStates[c.id] = PlayerState{Pos: pos, Energy: 10, Action: "", Advanced: false, Player: len(gs.PlayerStates) + 1}
	}

	emit(Event{
		Type:    EventJoinRoom,
		Payload: json.RawMessage(fmt.Sprintf(`{"room": "%s"}`, joinRoomEvent.Room)),
	}, c)

	// get list of rooms
	ListRoomHandler(Event{Type: EventListRooms}, c)

	status := "Waiting for opponent to arrive"
	if len(gs.PlayerStates) > 1 {
		status = "Game in progress, choose an action!"
	}
	roomEmit(Event{
		Type:    "UPDATE_STATUS",
		Payload: json.RawMessage(fmt.Sprintf(`{"status": "%s"}`, status)),
	}, c.room, c.hub)

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
	roomEmit(outgoingEvent, c.room, c.hub)

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

// Send an event to a single client
func emit(event Event, client *Client) {
	client.egress <- event
}

// Send an event to all connected clients
// func broadcast(event Event, h *Hub) {
// 	for client := range h.client {
// 		client.egress <- event
// 	}
// }

// Send an event to all clients in the same room
func roomEmit(event Event, room string, h *Hub) {
	for client := range h.client {
		if client.room == room {
			client.egress <- event
		}
	}
}
