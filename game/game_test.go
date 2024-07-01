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

func TestFirstMoveSuccess(t *testing.T) {
	g := game.New()

	nextPlayer, status, err := g.Sow(game.PLAYER_ONE, 2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if nextPlayer != game.PLAYER_ONE {
		t.Error("The next player should be player one, since they ended with a scoring move")
	}

	if status != game.STARTED {
		t.Error("The game status should be 'STARTED'")
	}

	state := g.GetState()

	expectedPits := [2][game.SIZE]uint{{4, 4, 0, 5, 5, 5}, {4, 4, 4, 4, 4, 4}}
	actualPits := state.Pits
	pitsEqual := true
	for i := range expectedPits {
		for j := range expectedPits[i] {
			// t.Logf("%d, %d", i, j)
			if expectedPits[i][j] != actualPits[i][j] {
				pitsEqual = false
			}
		}
	}

	if !pitsEqual {
		t.Errorf("Pits are not equal! Expected %v, got %v\n", expectedPits, actualPits)
	}

	expectedStores := [2]uint{1, 0}
	actualStores := state.Stores
	storesEqual := true
	for i := range expectedStores {
		if expectedStores[i] != actualStores[i] {
			storesEqual = false
		}
	}

	if !storesEqual {
		t.Errorf("Stores are not equal! Expected %v, got %v\n", expectedStores, actualStores)
	}
}
