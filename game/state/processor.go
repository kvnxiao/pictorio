package state

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/kvnxiao/pictorio/events"
	"github.com/kvnxiao/pictorio/game/state/chat"
	"github.com/kvnxiao/pictorio/game/user"
	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type GameState interface {
	EventProcessor(cleanupChan chan bool)
	Cleanup() <-chan bool

	// Status gets the GameStatus of the current game
	Status() model.GameStatus
	IsFull() bool

	StartGame()
	NextTurn()

	UserJoined(ctx context.Context, user *user.User, connErrChan chan error)
	UserLeft(userID string)
}

// GameStateProcessor handles the state of the game
type GameStateProcessor struct {
	mu sync.RWMutex

	maxPlayers int

	// status is the current GameStatus of the game
	status model.GameStatus

	// playerStates represents the userID -> player states mapping
	playerStates map[string]PlayerState

	// playerOrder represents the order for players (randomized on game start)
	playerOrder []string

	// The current word to guess
	currentWord string

	// currentTurn represents the player ID for the current player's turn
	currentTurn string

	messageQueue chan []byte

	chatHistory *chat.Chat

	// cleanedUpChan represents whether or not the game state has been cleaned up for the current room
	cleanedUpChan chan bool
}

func NewGameStateProcessor(maxPlayers int) GameState {
	return &GameStateProcessor{
		maxPlayers:    maxPlayers,
		status:        model.StatusWaitingReadyUp,
		playerStates:  make(map[string]PlayerState),
		playerOrder:   []string{},
		messageQueue:  make(chan []byte),
		chatHistory:   chat.NewChatHistory(),
		cleanedUpChan: make(chan bool),
	}
}

func (g *GameStateProcessor) playerStatesList() []model.PlayerState {
	g.mu.Lock()
	defer g.mu.Unlock()

	var playerStates []model.PlayerState

	for _, p := range g.playerStates {
		playerStates = append(playerStates, model.PlayerState{
			User:        p.UserModel(),
			Points:      p.Points(),
			Wins:        p.Wins(),
			IsSpectator: p.IsSpectator(),
			IsConnected: p.IsConnected(),
			IsReady:     p.IsReady(),
		})
	}

	return playerStates
}

// EventLoop represents the single-threaded game logic, which handles and processes incoming WebSocket messages from
// players, as well as handles cleaning up the room when all users have left the room.
func (g *GameStateProcessor) EventProcessor(cleanupChan chan bool) {
	for {
		select {
		case msg := <-g.messageQueue:
			var event events.GameEvent
			err := json.Unmarshal(msg, &event)
			if err != nil {
				log.Error().
					Bytes("msg", msg).
					Err(err).
					Msg("Failed to parse incoming user event")
			}

			switch event.Type {
			case events.EventTypeUserJoinLeaveEvent:
				g.onUserJoinLeaveEvent()
			case events.EventTypeRehydrate:
				g.onRehydrateEvent()
			case events.EventTypeChat:
				var chatEvent events.ChatEvent
				err := json.Unmarshal(event.Data, &chatEvent)
				if err != nil {
					log.Error().
						Err(err).Msg("Could not unmarshal ChatEvent from user")
				}
				g.onChatEvent(chatEvent)
			case events.EventTypeDraw:
				g.onDrawEvent(event)
			default:
				log.Error().Msg("Unknown event type unmarshalled from incoming user event")
			}
		case <-cleanupChan:
			g.cleanup()
			return
		}
	}
}

func (g *GameStateProcessor) cleanup() {
	log.Info().Msg("Cleaning up game state processor.")

	// TODO: cleanup game state processor
	g.chatHistory.Clear()

	log.Info().Msg("Done cleaning up game state processor!")

	// after cleanup, signal to the room that the game state is done cleaning up
	g.cleanedUpChan <- true
}

func (g *GameStateProcessor) Cleanup() <-chan bool {
	return g.cleanedUpChan
}

func (g *GameStateProcessor) Status() model.GameStatus {
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
