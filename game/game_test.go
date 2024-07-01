package game_test

import (
	"testing"

	"github.com/Samathingamajig/mancala/game"
)

func TestNewGameCreation(t *testing.T) {
	g := game.New() // implied: shouldn't panic
	if g == nil {
		t.Error("Constructing a new game shouldn't be nil")
	}
}

func TestNewGameState(t *testing.T) {
	g := game.New()

	state := g.GetState()

	// Pits
	hadIssuePits := false
	for i, row := range state.Pits {
		for j, seeds := range row {
			if state.Pits[i][j] != game.INITIAL_SEEDS {
				t.Errorf("The initial value for each pit should be %d, got %d at %d, %d\n", game.INITIAL_SEEDS, seeds, i, j)
				hadIssuePits = true
			}
		}
	}

	if hadIssuePits {
		t.Errorf("Got: %v\n", state.Pits)
	}

	// Stores
	for i, seeds := range state.Stores {
		if seeds != 0 {
			t.Errorf("The initial value for each store should be 0, got %d at index %d\n", seeds, i)
		}
	}

	// Player
	if state.NextPlayer != game.PLAYER_ONE {
		t.Error("The first player should be player one\n")
	}

	// Status
	if state.Status != game.FRESH {
		t.Error("A new game should be the 'FRESH' status\n")
	}
}
