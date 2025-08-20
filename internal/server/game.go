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

	if len(ids) != 2 {
		emit(Event{
			Type:    "WAITING_FOR_OPPONENT",
			Payload: json.RawMessage(`{"message": "Waiting for opponent to arrive"}`),
		}, c)
		return nil
	}
	// Game logic
	p1 := gs.PlayerStates[ids[0]]
	p2 := gs.PlayerStates[ids[1]]

	// Check if both players have made their actions
	if p1.Action == "" || p2.Action == "" {
		emit(Event{
			Type:    "WAITING_FOR_OPPONENT",
			Payload: json.RawMessage(`{"message": "Waiting for opponent to act"}`),
		}, c)
		return nil
	}

	// Simultaneous movement resolution
	var intendedPos1, intendedPos2 int
	var movementLog1, movementLog2 string
	intendedPos1, p1.Energy, movementLog1 = p1.resolveIntendedMovement()
	intendedPos2, p2.Energy, movementLog2 = p2.resolveIntendedMovement()
	p1.Pos, p2.Pos = resolveSimultaneousMovement(p1.Pos, intendedPos1, p2.Pos, intendedPos2)

	// Then resolve combat for both players
	var combatLog1, combatLog2 string
	p1.Energy, combatLog1 = p1.resolveCombat(&p2)
	p2.Energy, combatLog2 = p2.resolveCombat(&p1)

	// Update game state for next round
	gs.PlayerStates[ids[0]] = p1
	gs.PlayerStates[ids[1]] = p2
	gs.Turn++
	gs.LastAction = fmt.Sprintf("P1: %s %s| P2: %s %s.", movementLog1, combatLog1, movementLog2, combatLog2)

	// Send personalized state and winner to each client
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
			// Set personalized winner message
			if youState.Energy > 0 && opponentState.Energy <= 0 {
				personalized.Winner = "You win!"
			} else if youState.Energy <= 0 && opponentState.Energy > 0 {
				personalized.Winner = "Opponent wins!"
			} else if youState.Energy <= 0 && opponentState.Energy <= 0 {
				personalized.Winner = "Draw!"
			} else if gs.Turn > gs.MaxTurns {
				if youState.Energy > opponentState.Energy {
					personalized.Winner = "You win by energy!"
				} else if youState.Energy < opponentState.Energy {
					personalized.Winner = "Opponent wins by energy!"
				} else {
					personalized.Winner = "Draw!"
				}
			}

			// If there is a winner, set personalized.GameOver to true
			if personalized.Winner != "" {
				personalized.GameOver = true
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

// Movement actions: WAIT, RETREAT, ADVANCE
// Returns intended new position, new energy, and movement log
func (p1 *PlayerState) resolveIntendedMovement() (int, int, string) {
	switch p1.Action {
	case "WAIT":
		return p1.Pos, p1.Energy + 1, "WAIT: +1 Energy"
	case "RETREAT":
		var newPos int
		if p1.Player == 1 {
			newPos = max(0, p1.Pos-1)
		} else {
			newPos = min(6, p1.Pos+1)
		}
		return newPos, p1.Energy - 1, "RETREAT: -1 Energy"
	case "ADVANCE":
		var newPos int
		if p1.Player == 1 {
			newPos = min(6, p1.Pos+1)
		} else {
			newPos = max(0, p1.Pos-1)
		}
		p1.Advanced = true
		return newPos, p1.Energy - 1, "ADVANCE: Double attack next turn"
	}
	return p1.Pos, p1.Energy, ""
}

// Combat actions: ATTACK, COUNTER
func (p1 *PlayerState) resolveCombat(p2 *PlayerState) (int, string) {
	log := ""
	adjacent := abs(p1.Pos-p2.Pos) == 1

	switch p1.Action {
	case "ATTACK":
		dmg := 3
		if p1.Advanced {
			dmg = 6
		}
		switch p2.Action {
		case "COUNTER":
			if adjacent {
				p1.Energy -= dmg
				log = fmt.Sprintf("Opponent countered! takes %d damage.", dmg)
			} else {
				log = "Countered, but not adjacent. No damage."
				p1.Energy--
			}
		case "RETREAT":
			log = "Attack missed! Opponent retreated."
			p1.Energy--
		default:
			if adjacent {
				p2.Energy -= dmg
				log = fmt.Sprintf("Attacked for %d damage.", dmg)
			} else {
				log = "Attack missed! Not adjacent."
				p1.Energy--
			}
		}
	case "COUNTER":
		if p2.Action != "ATTACK" || !adjacent {
			p1.Energy -= 2
			log = "Countered nothing. -2 Energy."
		}
	}
	if p1.Action != "ADVANCE" {
		p1.Advanced = false
	}
	return p1.Energy, log
}

// resolveSimultaneousMovement applies blocking, swap, and priority rules and returns new positions for both players
func resolveSimultaneousMovement(pos1, intended1, pos2, intended2 int) (newPos1, newPos2 int) {
	// Swap places: both blocked
	if intended1 == pos2 && intended2 == pos1 {
		return pos1, pos2
	}
	// Both try to move to same square (not their current): player 1 gets priority
	if intended1 == intended2 && intended1 != pos1 && intended2 != pos2 {
		return intended1, pos2
	}
	// p1 blocked by p2
	if intended1 == pos2 && intended2 == pos2 {
		return pos1, intended2
	}
	// p2 blocked by p1
	if intended2 == pos1 && intended1 == pos1 {
		return intended1, pos2
	}
	// Otherwise, both move
	return intended1, intended2
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
