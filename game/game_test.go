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

func TestFirstMoveFailure(t *testing.T) {
	g := game.New()

	_, _, err := g.Sow(game.PLAYER_ONE, 6)
	if err == nil {
		t.Error("Expected an error when trying to sow from a store")
	}
}

func TestSingleMoveGames(t *testing.T) {
	cases := []struct {
		player   game.Player
		pitIndex uint
		expected game.MancalaGameState
	}{
		{
			game.PLAYER_ONE,
			0,
			game.MancalaGameState{
				Pits:       [2][game.SIZE]uint{{0, 5, 5, 5, 5, 4}, {4, 4, 4, 4, 4, 4}},
				Stores:     [2]uint{0, 0},
				NextPlayer: game.PLAYER_TWO,
				Status:     game.STARTED,
			},
		},
		{
			game.PLAYER_ONE,
			1,
			game.MancalaGameState{
				Pits:       [2][game.SIZE]uint{{4, 0, 5, 5, 5, 5}, {4, 4, 4, 4, 4, 4}},
				Stores:     [2]uint{0, 0},
				NextPlayer: game.PLAYER_TWO,
				Status:     game.STARTED,
			},
		},
		{
			game.PLAYER_ONE,
			2,
			game.MancalaGameState{
				Pits:       [2][game.SIZE]uint{{4, 4, 0, 5, 5, 5}, {4, 4, 4, 4, 4, 4}},
				Stores:     [2]uint{1, 0},
				NextPlayer: game.PLAYER_ONE,
				Status:     game.STARTED,
			},
		},
		{
			game.PLAYER_ONE,
			3,
			game.MancalaGameState{
				Pits:       [2][game.SIZE]uint{{4, 4, 4, 0, 5, 5}, {5, 4, 4, 4, 4, 4}},
				Stores:     [2]uint{1, 0},
				NextPlayer: game.PLAYER_TWO,
				Status:     game.STARTED,
			},
		},
		{
			game.PLAYER_ONE,
			4,
			game.MancalaGameState{
				Pits:       [2][game.SIZE]uint{{4, 4, 4, 4, 0, 5}, {5, 5, 4, 4, 4, 4}},
				Stores:     [2]uint{1, 0},
				NextPlayer: game.PLAYER_TWO,
				Status:     game.STARTED,
			},
		},
		{
			game.PLAYER_ONE,
			5,
			game.MancalaGameState{
				Pits:       [2][game.SIZE]uint{{4, 4, 4, 4, 4, 0}, {5, 5, 5, 4, 4, 4}},
				Stores:     [2]uint{1, 0},
				NextPlayer: game.PLAYER_TWO,
				Status:     game.STARTED,
			},
		},
	}

	for caseIdx, c := range cases {
		state, err := playGame([]struct {
			player   game.Player
			position uint
		}{{c.player, c.pitIndex}})

		if err != nil {
			t.Errorf("Error playing game %d: %v", caseIdx, err)
			t.FailNow()
		}

		if !isEqualGameState(c.expected, state) {
			t.Errorf("Game %d failed! Expected %v, got %v", caseIdx, c.expected, state)
		}
	}
}

func TestMultipleMoveGameWithCapture(t *testing.T) {
	expectedState := game.MancalaGameState{
		Pits:       [2][game.SIZE]uint{{0, 1, 3, 8, 7, 0}, {0, 0, 6, 6, 5, 0}},
		Stores:     [2]uint{10, 2},
		NextPlayer: game.PLAYER_TWO,
		Status:     game.STARTED,
	}

	state, err := playGame([]struct {
		player   game.Player
		position uint
	}{
		{game.PLAYER_ONE, 2},
		{game.PLAYER_ONE, 5},
		{game.PLAYER_TWO, 1},
		{game.PLAYER_TWO, 5},
		{game.PLAYER_ONE, 1},
		{game.PLAYER_ONE, 5},
		{game.PLAYER_ONE, 0},
	})

	if err != nil {
		t.Errorf("Error playing game: %v", err)
		t.FailNow()
	}

	if !isEqualGameState(expectedState, state) {
		t.Errorf("Game failed! Expected %v, got %v", expectedState, state)
	}
}

func playGame(moves []struct {
	player   game.Player
	position uint
}) (game.MancalaGameState, error) {
	g := game.New()

	for _, move := range moves {
		_, _, err := g.Sow(move.player, move.position)
		if err != nil {
			return g.GetState(), err
		}
	}

	return g.GetState(), nil
}

func isEqualGameState(expected game.MancalaGameState, actual game.MancalaGameState) bool {
	expectedPits := expected.Pits
	actualPits := actual.Pits
	for i := range expectedPits {
		for j := range expectedPits[i] {
			if expectedPits[i][j] != actualPits[i][j] {
				return false
			}
		}
	}

	expectedStores := expected.Stores
	actualStores := actual.Stores
	for i := range expectedStores {
		if expectedStores[i] != actualStores[i] {
			return false
		}
	}

	return true
}
