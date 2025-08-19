package server

import (
	"encoding/json"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
	EventListRooms   = "list_rooms"
	EventJoinRoom    = "join_room"
	EventLeaveRoom   = "leave_room"
	EventInitClient  = "init_client"
	EventGameAction  = "game_action"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type JoinRoomEvent struct {
	Room string `json:"room"`
}

type ListRoomEvent struct {
	// No fields needed for this event
}

type ListRoomResponse struct {
	Rooms []string `json:"rooms"`
}

type InitClientEvent struct {
	PlayerId string `json:"playerId"`
}
