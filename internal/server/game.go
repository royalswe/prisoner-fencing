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
			Type: "WAITING_FOR_OPPONENT",
		}, c)
		return nil
	}

	// First resolve movement for both players
	var movementLog1, movementLog2 string
	p1.Pos, p1.Energy, movementLog1 = p1.resolveMovement(&p2)
	p2.Pos, p2.Energy, movementLog2 = p2.resolveMovement(&p1)
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
func (p1 *PlayerState) resolveMovement(p2 *PlayerState) (int, int, string) {
	log := ""
	switch p1.Action {
	case "WAIT":
		p1.Energy++
		log = "WAIT: +1 Energy"
	case "RETREAT":
		var newPos int
		if p1.Player == 1 {
			newPos = max(0, p1.Pos-1)
		} else {
			newPos = min(6, p1.Pos+1)
		}
		if newPos != p2.Pos {
			p1.Pos = newPos
		}
		p1.Energy--
		log = "RETREAT: -1 Energy"
	case "ADVANCE":
		var newPos int
		if p1.Player == 1 {
			newPos = min(6, p1.Pos+1)
		} else {
			newPos = max(0, p1.Pos-1)
		}
		if newPos != p2.Pos {
			p1.Pos = newPos
		}
		p1.Energy--
		p1.Advanced = true
		log = "ADVANCE: Double attack next turn"
	}
	return p1.Pos, p1.Energy, log
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

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
