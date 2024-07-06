package game

import (
	"fmt"
)

const SIZE = 6
const INITIAL_SEEDS = 4

type Player uint8
type Status uint8

const (
	PLAYER_ONE Player = iota
	PLAYER_TWO
)

const (
	FRESH Status = iota
	STARTED
	FINISHED
)

type MancalaGameState struct {
	Pits       [2][SIZE]uint
	Stores     [2]uint
	NextPlayer Player
	Status     Status
}

type MancalaGame struct {
	pits       [SIZE*2 + 2]uint
	nextPlayer Player
	status     Status
}

func New() *MancalaGame {
	game := &MancalaGame{
		pits:       [SIZE*2 + 2]uint{},
		nextPlayer: PLAYER_ONE,
		status:     FRESH,
	}
	for i := 0; i < SIZE; i++ {
		game.pits[i] = INITIAL_SEEDS
		game.pits[i+SIZE+1] = INITIAL_SEEDS
	}

	return game
}

// Returns the current state of the Mancala game
// Clients are responsible for any rotating of boards
// (i.e. pits are returned from the POV of each respective player viewing)
func (g *MancalaGame) GetState() MancalaGameState {
	state := MancalaGameState{
		Pits:       [2][SIZE]uint{},
		Stores:     [2]uint{},
		NextPlayer: g.nextPlayer,
		Status:     g.status,
	}

	for i := 0; i < SIZE; i++ {
		state.Pits[0][i] = g.pits[i]
		state.Pits[1][i] = g.pits[SIZE+1+i]
	}

	state.Stores[0] = g.pits[SIZE]
	state.Stores[1] = g.pits[SIZE*2+1]

	return state
}

func (g *MancalaGame) Sow(player Player, position uint) (Player, Status, error) {
	if g.status == FINISHED {
		return g.nextPlayer, g.status, fmt.Errorf("attempted to sow after game is over")
	}

	if g.nextPlayer != player {
		return g.nextPlayer, g.status, fmt.Errorf("attempted to sow with the wrong next player")
	}

	if position >= SIZE {
		return g.nextPlayer, g.status, fmt.Errorf("attempted to sow out of range: %d", position)
	}

	idx := position
	if player == PLAYER_TWO {
		idx += SIZE + 1
	}

	seeds := g.pits[idx]
	if seeds == 0 {
		return g.nextPlayer, g.status, fmt.Errorf("attempted to sow from an empty pit")
	}
	g.pits[idx] = 0

	g.status = STARTED

	for ; seeds > 0; seeds-- {
		idx++
		if (idx == SIZE && player == PLAYER_TWO) || (idx == 2*SIZE+1 && player == PLAYER_ONE) {
			// In a store not belonging to a player, proceed to next pit
			idx++
		}
		if idx == SIZE*2+2 {
			// Wrap around
			idx = 0
		}

		g.pits[idx]++
	}

	if g.pits[idx] == 1 && idx != SIZE && idx != SIZE*2+1 {
		// Try to capture seeds
		otherIdx := 2*SIZE - idx
		var storeIdx uint
		canCapture := false

		if player == PLAYER_ONE && idx < SIZE {
			storeIdx = SIZE
			canCapture = true
		} else if player == PLAYER_TWO && idx > SIZE+1 {
			otherIdx = 2*SIZE - idx
			storeIdx = SIZE*2 + 1
			canCapture = true
		}

		if canCapture && g.pits[otherIdx] > 0 {
			g.pits[storeIdx] += g.pits[idx] + g.pits[otherIdx]
			g.pits[idx] = 0
			g.pits[otherIdx] = 0
		}
	}

	// Swap players, unless ending in a store
	if idx%(SIZE+1) != SIZE {
		if player == PLAYER_ONE {
			g.nextPlayer = PLAYER_TWO
		} else {
			g.nextPlayer = PLAYER_ONE
		}
	}

	// Ending states (empty row for either player, player with seeds on their side can collect for their store)
	playerOneSum := uint(0)
	playerTwoSum := uint(0)
	for i := 0; i < SIZE; i++ {
		playerOneSum += g.pits[i]
		playerTwoSum += g.pits[SIZE+1+i]
	}

	if playerOneSum == 0 || playerTwoSum == 0 {
		for i := 0; i < SIZE; i++ {
			g.pits[i] = 0
			g.pits[SIZE+1+i] = 0
		}
		g.pits[SIZE] += playerOneSum
		g.pits[SIZE*2+1] += playerTwoSum
		g.status = FINISHED
	}

	return g.nextPlayer, g.status, nil
}
