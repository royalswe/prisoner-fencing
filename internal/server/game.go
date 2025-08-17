package server

import (
	"encoding/json"
	"fmt"
)

type GameState struct {
	Turn         int                    `json:"turn"`
	MaxTurns     int                    `json:"maxTurns"`
	LastAction   string                 `json:"lastAction"`
	GameOver     bool                   `json:"gameOver"`
	Winner       string                 `json:"winner"`
	PlayerStates map[string]PlayerState `json:"playerStates"` // id -> state
}

type PlayerState struct {
	Pos      int    `json:"pos"`
	Energy   int    `json:"energy"`
	Action   string `json:"action"`
	Advanced bool   `json:"advanced"`
	Player   int    `json:"player"` // 1 for you, 2 for opponent
}

var roomStates = make(map[string]*GameState)

func GameActionHandler(event Event, c *Client) error {

	var payload struct {
		Room     string `json:"room"`
		PlayerId string `json:"playerId"`
		Action   string `json:"action"`
	}
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal game action: %v", err)
	}

	gs, ok := roomStates[payload.Room]
	if !ok {
		gs = &GameState{
			Turn:         1,
			MaxTurns:     20,
			PlayerStates: make(map[string]PlayerState),
		}
		roomStates[payload.Room] = gs
	}

	// Set or update player state
	if _, exists := gs.PlayerStates[c.id]; !exists {
		pos := 2
		if len(gs.PlayerStates) > 0 {
			pos = 4
		}
		gs.PlayerStates[c.id] = PlayerState{Pos: pos, Energy: 10, Action: payload.Action, Advanced: false, Player: len(gs.PlayerStates) + 1}
	} else {
		ps := gs.PlayerStates[c.id]
		ps.Action = payload.Action
		gs.PlayerStates[c.id] = ps
	}

	// Get both player ids
	var ids []string
	for pid := range gs.PlayerStates {
		ids = append(ids, pid)
	}
	fmt.Print(ids)

	if len(ids) != 2 {
		fmt.Printf("expected 2 players, got %d", len(ids))
		fmt.Printf("ids: %+v\n", ids)
		return nil
	}
	// Game logic
	p1 := gs.PlayerStates[ids[0]]
	p2 := gs.PlayerStates[ids[1]]

	// Check if both players have made their actions
	if p1.Action == "" || p2.Action == "" {
		// Not enough players, wait for the opponent
		emit(Event{
			Type:    "WAITING_FOR_OPPONENT",
			Payload: json.RawMessage(fmt.Sprintf(`{"waitingForOpponent": true, "room": "%s"}`, payload.Room)),
		}, c)
		return nil
	}

	// For player 2, invert movement direction
	var log1, log2 string
	p1.Pos, p1.Energy, log1 = p1.resolveAction(&p2, false)
	p2.Pos, p2.Energy, log2 = p2.resolveAction(&p1, true)

	// ATTACK and COUNTER logic
	lastAction := ""

	// Update game state for next round
	gs.PlayerStates[ids[0]] = p1
	gs.PlayerStates[ids[1]] = p2
	gs.Turn++
	gs.LastAction = fmt.Sprintf("P1: %s, P2: %s. %s", log1, log2, lastAction)

	// Win conditions
	gameOver := false
	winner := ""
	if p1.Energy <= 0 {
		gameOver = true
		winner = "Opponent wins!"
	} else if p2.Energy <= 0 {
		gameOver = true
		winner = "You win!"
	} else if gs.Turn > gs.MaxTurns {
		gameOver = true
		if p1.Energy > p2.Energy {
			winner = "You win by energy!"
		} else if p1.Energy < p2.Energy {
			winner = "Opponent wins by energy!"
		} else {
			winner = "Draw!"
		}
	}
	gs.GameOver = gameOver
	gs.Winner = winner

	// Send personalized state to each client
	for client := range c.hub.client {
		if client.room == payload.Room {
			// Assign 'you' and 'opponent' based on client.id
			var youState, opponentState PlayerState
			for pid, state := range gs.PlayerStates {
				if pid == client.id {
					youState = state
				} else {
					opponentState = state
				}
			}
			personalized := *gs
			personalized.PlayerStates = map[string]PlayerState{
				"you":      youState,
				"opponent": opponentState,
			}
			data, _ := json.Marshal(personalized)
			outgoing := Event{
				Type:    "GAME_ACTION_RESULT",
				Payload: data,
			}
			emit(outgoing, client)
		}
	}

	p1.Action = ""
	p2.Action = ""
	gs.PlayerStates[ids[0]] = p1
	gs.PlayerStates[ids[1]] = p2
	return nil
}

// Helper to resolve actions
func (p1 *PlayerState) resolveAction(p2 *PlayerState, invert bool) (int, int, string) {
	log := ""
	switch p1.Action {
	case "WAIT":
		p1.Energy++
		log = "WAIT: +1 Energy"
	case "RETREAT":
		if invert {
			p1.Pos = min(6, p1.Pos+1)
		} else {
			p1.Pos = max(0, p1.Pos-1)
		}
		p1.Energy--
		log = "RETREAT: -1 Energy"
	case "ADVANCE":
		if invert {
			p1.Pos = max(0, p1.Pos-1)
			if p1.Pos <= p2.Pos {
				// can not advance through or same position as opponent
				p1.Pos = max(0, p1.Pos+1)
			}
		} else {
			p1.Pos = min(6, p1.Pos+1)
			if p1.Pos >= p2.Pos {
				p1.Pos = min(6, p1.Pos-1)
			}
		}
		p1.Energy--
		p1.Advanced = true
		log = "ADVANCE: Double attack next turn"
	case "ATTACK":
		dmg := 3
		if p1.Advanced {
			dmg = 6
		}
		if p2.Action == "COUNTER" {
			p1.Energy -= dmg
			log = fmt.Sprintf("Opponent countered! takes %d damage.", dmg)
		} else {
			p2.Energy -= dmg
			log = fmt.Sprintf("Attacked for %d damage.", dmg)
		}
	case "COUNTER":
		if p2.Action != "ATTACK" {
			p1.Energy -= 2
			log = "Countered nothing. -2 Energy."
		}
	}
	if p1.Action != "ADVANCE" {
		p1.Advanced = false
	}
	return p1.Pos, p1.Energy, log
}
