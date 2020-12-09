package state

import (
	"context"
	"sync"

	"github.com/kvnxiao/pictorio/game/user"
	"github.com/rs/zerolog/log"
)

type GameStatus int

const (
	StatusWaitingReadyUp GameStatus = iota
	StatusStarted
	StatusGameOver
)

type GameState interface {
	EventProcessor(cleanupChan chan bool)
	Cleanup() <-chan bool

	// Status gets the GameStatus of the current game
	Status() GameStatus
	IsFull() bool

	StartGame()
	NextTurn()

	UserJoined(ctx context.Context, user *user.User, connErrChan chan error)
	UserLeft(userID string)
}

// GameStateProcessor handles the state of the game
type GameStateProcessor struct {
	mu sync.Mutex

	maxPlayers int

	// status is the current GameStatus of the game
	status GameStatus

	// The current word to guess
	currentWord string

	// playerStates represents the player states
	playerStates map[string]PlayerState

	// playerOrder represents the order for players (randomized on game start)
	playerOrder []string

	// currentTurn represents the player ID for the current player's turn
	currentTurn string

	messageQueue chan []byte

	// cleanedUpChan represents whether or not the game state has been cleaned up for the current room
	cleanedUpChan chan bool
}

func NewGameStateProcessor(maxPlayers int) GameState {
	return &GameStateProcessor{
		maxPlayers:    maxPlayers,
		status:        StatusWaitingReadyUp,
		playerStates:  make(map[string]PlayerState),
		playerOrder:   []string{},
		messageQueue:  make(chan []byte),
		cleanedUpChan: make(chan bool),
	}
}

// EventLoop represents the single-threaded game logic, which handles and processes incoming WebSocket messages from
// players, as well as handles cleaning up the room when all users have left the room.
func (g *GameStateProcessor) EventProcessor(cleanupChan chan bool) {
	for {
		select {
		case _ = <-g.messageQueue:
			// TODO: unmarshal message as event, send to event handler
		case <-cleanupChan:
			g.cleanup()
			return
		}
	}
}

func (g *GameStateProcessor) cleanup() {
	log.Info().Msg("Cleaning up game state processor.")

	// TODO: cleanup game state processor

	log.Info().Msg("Done cleaning up game state processor!")

	// after cleanup, signal to the room that the game state is done cleaning up
	g.cleanedUpChan <- true
}

func (g *GameStateProcessor) Cleanup() <-chan bool {
	return g.cleanedUpChan
}

func (g *GameStateProcessor) Status() GameStatus {
	return g.status
}

func (g *GameStateProcessor) StartGame() {
	panic("implement me")
	// TODO: check all players are ready
	// TODO: randomize player turn order
	// TODO: progress in game event loop
}

func (g *GameStateProcessor) NextTurn() {
	panic("implement me")
}

func (g *GameStateProcessor) IsFull() bool {
	return len(g.playerStates) >= g.maxPlayers
}

func (g *GameStateProcessor) UserJoined(ctx context.Context, user *user.User, connErrChan chan error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	playerState, ok := g.playerStates[user.ID]
	if ok {
		// Existing user re-joined the room
		playerState.SetNewConnection(user)
	} else {
		// No player state exists for this user, create a new player state for the current game
		playerState = NewPlayer(user, g.IsFull())
		g.playerStates[user.ID] = playerState
	}
	playerState.SetConnected(true)

	go user.ReaderLoop(ctx, g.messageQueue, connErrChan)
	go user.WriterLoop(ctx, connErrChan)

	// TODO: send rehydration event to user who just joined
	// TODO: broadcast player join event
}

func (g *GameStateProcessor) UserLeft(userID string) {
	g.mu.Lock()
	defer g.mu.Unlock()

	playerState, ok := g.playerStates[userID]
	if !ok {
		return
	}

	playerState.SetConnected(false)
	// TODO: broadcast
}
