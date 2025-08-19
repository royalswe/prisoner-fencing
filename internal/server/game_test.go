package server

import (
	"testing"
)

func TestAdvanceDoubleAttack(t *testing.T) {
	gs := &GameState{
		Turn:         1,
		MaxTurns:     20,
		PlayerStates: map[string]PlayerState{},
	}
	gs.PlayerStates["p1"] = PlayerState{Pos: 2, Energy: 10, Action: "ADVANCE", Advanced: false, Player: 1}
	gs.PlayerStates["p2"] = PlayerState{Pos: 4, Energy: 10, Action: "WAIT", Advanced: false, Player: 2}

	p1 := gs.PlayerStates["p1"]
	p2 := gs.PlayerStates["p2"]
	p1.Action = "ADVANCE"
	p2.Action = "ADVANCE"
	p1.Pos, p1.Energy, _ = p1.resolveMovement(&p2)
	p2.Pos, p2.Energy, _ = p2.resolveMovement(&p1)
	p1.Energy, _ = p1.resolveCombat(&p2)
	p2.Energy, _ = p2.resolveCombat(&p1)
	gs.PlayerStates["p1"] = p1
	gs.PlayerStates["p2"] = p2

	if p1.Energy != 9 {
		t.Errorf("Expected p1 energy to be 9 after double attack, got %d", p1.Energy)
	}
	// Next turn: p1 attacks, p2 waits
	p1.Action = "ATTACK"
	p2.Action = "WAIT"
	p1.Pos, p1.Energy, _ = p1.resolveMovement(&p2)
	p2.Pos, p2.Energy, _ = p2.resolveMovement(&p1)
	p1.Energy, _ = p1.resolveCombat(&p2)
	p2.Energy, _ = p2.resolveCombat(&p1)

	if p2.Energy != 4 {
		t.Errorf("Expected p2 energy to be 4 after double attack, got %d", p2.Energy)
	}
	if p1.Advanced {
		t.Errorf("Advanced should reset after attack")
	}
}

func TestCounterReflect(t *testing.T) {
	gs := &GameState{
		Turn:         1,
		MaxTurns:     20,
		PlayerStates: map[string]PlayerState{},
	}
	gs.PlayerStates["p1"] = PlayerState{Pos: 3, Energy: 10, Action: "ATTACK", Advanced: false, Player: 1}
	gs.PlayerStates["p2"] = PlayerState{Pos: 4, Energy: 10, Action: "COUNTER", Advanced: false, Player: 2}

	p1 := gs.PlayerStates["p1"]
	p2 := gs.PlayerStates["p2"]
	p1.Pos, p1.Energy, _ = p1.resolveMovement(&p2)
	p2.Pos, p2.Energy, _ = p2.resolveMovement(&p1)
	p1.Energy, _ = p1.resolveCombat(&p2)
	p2.Energy, _ = p2.resolveCombat(&p1)

	if p1.Energy != 7 {
		t.Errorf("Expected p1 energy to be 7 after counter reflect, got %d", p1.Energy)
	}
	if p2.Energy != 10 {
		t.Errorf("Expected p2 energy to be 10 after counter reflect, got %d", p2.Energy)
	}
}

func TestCounterPenalty(t *testing.T) {
	gs := &GameState{
		Turn:         1,
		MaxTurns:     20,
		PlayerStates: map[string]PlayerState{},
	}
	gs.PlayerStates["p1"] = PlayerState{Pos: 2, Energy: 10, Action: "COUNTER", Advanced: false, Player: 1}
	gs.PlayerStates["p2"] = PlayerState{Pos: 4, Energy: 10, Action: "WAIT", Advanced: false, Player: 2}

	p1 := gs.PlayerStates["p1"]
	p2 := gs.PlayerStates["p2"]
	p1.Pos, p1.Energy, _ = p1.resolveMovement(&p2)
	p2.Pos, p2.Energy, _ = p2.resolveMovement(&p1)
	p1.Energy, _ = p1.resolveCombat(&p2)
	p2.Energy, _ = p2.resolveCombat(&p1)

	if p1.Energy != 8 {
		t.Errorf("Expected p1 energy to be 8 after counter penalty, got %d", p1.Energy)
	}
}

func TestWaitEnergy(t *testing.T) {
	gs := &GameState{
		Turn:         1,
		MaxTurns:     20,
		PlayerStates: map[string]PlayerState{},
	}
	gs.PlayerStates["p1"] = PlayerState{Pos: 2, Energy: 10, Action: "WAIT", Advanced: false, Player: 1}
	gs.PlayerStates["p2"] = PlayerState{Pos: 4, Energy: 10, Action: "WAIT", Advanced: false, Player: 2}

	p1 := gs.PlayerStates["p1"]
	p2 := gs.PlayerStates["p2"]

	p1.Pos, p1.Energy, _ = p1.resolveMovement(&p2)
	p2.Pos, p2.Energy, _ = p2.resolveMovement(&p1)
	p1.Pos, p1.Energy, _ = p1.resolveMovement(&p2)

	if p1.Energy != 12 {
		t.Errorf("Expected p1 energy to be 12 after WAIT twice, got %d", p1.Energy)
	}
	if p2.Energy != 11 {
		t.Errorf("Expected p2 energy to be 11 after WAIT, got %d", p2.Energy)
	}
}

func TestRetreatAdvanceEnergy(t *testing.T) {
	gs := &GameState{
		Turn:         1,
		MaxTurns:     20,
		PlayerStates: map[string]PlayerState{},
	}
	gs.PlayerStates["p1"] = PlayerState{Pos: 2, Energy: 10, Action: "RETREAT", Advanced: false, Player: 1}
	gs.PlayerStates["p2"] = PlayerState{Pos: 4, Energy: 10, Action: "ADVANCE", Advanced: false, Player: 2}

	p1 := gs.PlayerStates["p1"]
	p2 := gs.PlayerStates["p2"]
	p1.Pos, p1.Energy, _ = p1.resolveMovement(&p2)
	p2.Pos, p2.Energy, _ = p2.resolveMovement(&p1)

	if p1.Energy != 9 {
		t.Errorf("Expected p1 energy to be 9 after RETREAT, got %d", p1.Energy)
	}
	if p2.Energy != 9 {
		t.Errorf("Expected p2 energy to be 9 after ADVANCE, got %d", p2.Energy)
	}
	if !p2.Advanced {
		t.Errorf("Expected p2.Advanced to be true after ADVANCE")
	}
}

func TestAttackAdjacentHit(t *testing.T) {
	gs := &GameState{
		Turn:         1,
		MaxTurns:     20,
		PlayerStates: map[string]PlayerState{},
	}
	gs.PlayerStates["p1"] = PlayerState{Pos: 2, Energy: 10, Action: "ATTACK", Advanced: false, Player: 1}
	gs.PlayerStates["p2"] = PlayerState{Pos: 3, Energy: 10, Action: "WAIT", Advanced: false, Player: 2}

	p1 := gs.PlayerStates["p1"]
	p2 := gs.PlayerStates["p2"]
	p1.Pos, p1.Energy, _ = p1.resolveMovement(&p2)
	p2.Pos, p2.Energy, _ = p2.resolveMovement(&p1)
	p1.Energy, _ = p1.resolveCombat(&p2)
	p2.Energy, _ = p2.resolveCombat(&p1)

	if p2.Energy != 8 {
		t.Errorf("Expected p2 energy to be 8 after adjacent attack, got %d", p2.Energy)
	}
}

func TestAttackNotAdjacentMiss(t *testing.T) {
	gs := &GameState{
		Turn:         1,
		MaxTurns:     20,
		PlayerStates: map[string]PlayerState{},
	}
	gs.PlayerStates["p1"] = PlayerState{Pos: 2, Energy: 10, Action: "ATTACK", Advanced: false, Player: 1}
	gs.PlayerStates["p2"] = PlayerState{Pos: 4, Energy: 10, Action: "WAIT", Advanced: false, Player: 2}

	p1 := gs.PlayerStates["p1"]
	p2 := gs.PlayerStates["p2"]
	p1.Pos, p1.Energy, _ = p1.resolveMovement(&p2)
	p2.Pos, p2.Energy, _ = p2.resolveMovement(&p1)
	p1.Energy, _ = p1.resolveCombat(&p2)
	p2.Energy, _ = p2.resolveCombat(&p1)

	if p2.Energy != 11 {
		t.Errorf("Expected p2 energy to be 11 after non-adjacent attack, got %d", p2.Energy)
	}
	if p1.Energy != 9 {
		t.Errorf("Expected p1 energy to be 9 after missed attack, got %d", p1.Energy)
	}
}

func TestAttackMissOnRetreat(t *testing.T) {
	gs := &GameState{
		Turn:         1,
		MaxTurns:     20,
		PlayerStates: map[string]PlayerState{},
	}
	gs.PlayerStates["p1"] = PlayerState{Pos: 2, Energy: 10, Action: "ATTACK", Advanced: false, Player: 1}
	gs.PlayerStates["p2"] = PlayerState{Pos: 3, Energy: 10, Action: "RETREAT", Advanced: false, Player: 2}

	p1 := gs.PlayerStates["p1"]
	p2 := gs.PlayerStates["p2"]
	p1.Pos, p1.Energy, _ = p1.resolveMovement(&p2)
	p2.Pos, p2.Energy, _ = p2.resolveMovement(&p1)
	p1.Energy, _ = p1.resolveCombat(&p2)
	p2.Energy, _ = p2.resolveCombat(&p1)

	if p2.Energy != 9 {
		t.Errorf("Expected p2 energy to be 9 after retreat, got %d", p2.Energy)
	}
	if p1.Energy != 9 {
		t.Errorf("Expected p1 energy to be 9 after missed attack on retreat, got %d", p1.Energy)
	}
}

func TestCounterOnlyIfAdjacent(t *testing.T) {
	gs := &GameState{
		Turn:         1,
		MaxTurns:     20,
		PlayerStates: map[string]PlayerState{},
	}
	gs.PlayerStates["p1"] = PlayerState{Pos: 2, Energy: 10, Action: "ATTACK", Advanced: false, Player: 1}
	gs.PlayerStates["p2"] = PlayerState{Pos: 4, Energy: 10, Action: "COUNTER", Advanced: false, Player: 2}

	p1 := gs.PlayerStates["p1"]
	p2 := gs.PlayerStates["p2"]
	p1.Pos, p1.Energy, _ = p1.resolveMovement(&p2)
	p2.Pos, p2.Energy, _ = p2.resolveMovement(&p1)
	p1.Energy, _ = p1.resolveCombat(&p2)
	p2.Energy, _ = p2.resolveCombat(&p1)

	if p1.Energy != 9 {
		t.Errorf("Expected p1 energy to be 9 after non-adjacent counter, got %d", p1.Energy)
	}
	if p2.Energy != 8 {
		t.Errorf("Expected p2 energy to be 8 after non-adjacent counter, got %d", p2.Energy)
	}
}
