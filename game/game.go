package game

const SIZE = 6
const INITIAL_SEEDS = 4

type Player = uint8
type GameStatus = uint8

const (
	PLAYER_ONE Player = iota
	PLAYER_TWO
)

const (
	FRESH GameStatus = iota
	STARTED
	FINISHED
)

type MancalaGameState struct {
	Pits       [2][SIZE]uint
	Stores     [2]uint
	NextPlayer Player
	Status     GameStatus
}

type MancalaGame struct {
	Pits       [SIZE*2 + 2]uint
	nextPlayer Player
	status     GameStatus
}

func New() *MancalaGame {
	game := &MancalaGame{
		Pits:       [SIZE*2 + 2]uint{},
		nextPlayer: PLAYER_ONE,
		status:     FRESH,
	}
	for i := 0; i < SIZE; i++ {
		game.Pits[i] = INITIAL_SEEDS
		game.Pits[i+SIZE+1] = INITIAL_SEEDS
	}

	return game
}

func (g *MancalaGame) GetState() MancalaGameState {
	state := MancalaGameState{
		Pits:       [2][SIZE]uint{},
		Stores:     [2]uint{},
		NextPlayer: g.nextPlayer,
		Status:     g.status,
	}

	for i := 0; i < SIZE; i++ {
		state.Pits[0][i] = g.Pits[i]
		state.Pits[1][SIZE-1-i] = g.Pits[SIZE+1+i]
	}

	state.Stores[0] = g.Pits[SIZE]
	state.Stores[1] = g.Pits[SIZE*2+1]

	return state
}
