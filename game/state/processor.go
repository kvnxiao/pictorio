package state

import (
	"context"
	"encoding/json"
	"math/rand"

	"github.com/kvnxiao/pictorio/events"
	"github.com/kvnxiao/pictorio/game/state/chat"
	"github.com/kvnxiao/pictorio/game/state/drawing"
	"github.com/kvnxiao/pictorio/game/state/players"
	"github.com/kvnxiao/pictorio/game/state/status"
	"github.com/kvnxiao/pictorio/game/user"
	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type GameState interface {
	EventProcessor(cleanupChan chan bool)
	Cleanup() <-chan bool

	// Status gets the Status of the current game
	Status() model.GameStatus

	StartGame() bool
	NextTurn()

	HandleUserConnection(ctx context.Context, user *user.User, connErrChan chan error)
	RemoveUserConnection(userID string)
}

// GameStateProcessor handles the state of the game
type GameStateProcessor struct {
	// status is the current Status of the game
	status status.GameStatus

	// drawing is the current drawing history / state of the game
	drawing drawing.History

	// players represents the userID -> player states mapping
	players players.Players

	// chatHistory is the chat history since the beginning of the game
	chatHistory chat.History

	// cleanedUpChan represents whether or not the game state has been cleaned up for the current room
	cleanedUpChan chan bool

	// messageQueue represents the message queue for events incoming from individual user WebSockets, which will be
	// processed by the EventProcessor method to handle events
	messageQueue chan []byte

	// wordSelectionIndex allows the current turn player to send a TurnSelectionEvent which will be recorded by the
	// game processor
	wordSelectionIndex chan SelectionIndex

	// wordGuess allows the message queue to process a ChatEvent as a word guess when the model.GameStatus is started
	// and the model.TurnStatus is drawing
	wordGuess chan Guess
}

func NewGameStateProcessor(maxPlayers int) GameState {
	return &GameStateProcessor{
		status:             status.NewGameStatus(),
		drawing:            drawing.NewDrawingHistory(),
		players:            players.NewPlayerContainer(maxPlayers),
		chatHistory:        chat.NewChatHistory(),
		cleanedUpChan:      make(chan bool),
		messageQueue:       make(chan []byte),
		wordSelectionIndex: make(chan SelectionIndex),
		wordGuess:          make(chan Guess),
	}
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
			case events.EventTypeUserJoinLeave:
				g.onUserJoinLeaveEvent()

			case events.EventTypeRehydrate:
				g.onRehydrateEvent()

			case events.EventTypeChat:
				var chatEvent events.ChatEvent
				err := json.Unmarshal(event.Data, &chatEvent)
				if err != nil {
					log.Error().Err(err).
						Msg("Could not unmarshal " + events.EventTypeChat.String() + " from user")
				}
				g.onChatEvent(chatEvent)

			case events.EventTypeDraw:
				var drawEvent events.DrawEvent
				err := json.Unmarshal(event.Data, &drawEvent)
				if err != nil {
					log.Error().Err(err).
						Msg("Could not unmarshal " + events.EventTypeDraw.String() + " from user")
				}
				g.onDrawEvent(drawEvent)

			case events.EventTypeReady:
				var readyEvent events.ReadyEvent
				err := json.Unmarshal(event.Data, &readyEvent)
				if err != nil {
					log.Error().Err(err).
						Msg("Could not unmarshal " + events.EventTypeReady.String() + " from user")
				}
				g.onReadyEvent(readyEvent)

			case events.EventTypeStartGame:
				g.onStartGameEvent()

			case events.EventTypeStartGameIssued:
				var startGameIssuedEvent events.StartGameIssuedEvent
				err := json.Unmarshal(event.Data, &startGameIssuedEvent)
				if err != nil {
					log.Error().Err(err).
						Msg("Could not unmarshal " + events.EventTypeStartGameIssued.String() + " from user")
				}
				g.onStartGameIssuedEvent(startGameIssuedEvent)

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

	// Cleanup game state processor
	g.chatHistory.Clear()
	g.drawing.Clear()
	g.status.Cleanup()
	g.players.Cleanup()
	g.chatHistory = nil
	g.drawing = nil
	g.status = nil
	g.players = nil
	close(g.messageQueue)
	close(g.wordSelectionIndex)
	close(g.wordGuess)

	// TODO: close channels if necessary

	log.Info().Msg("Done cleaning up game state processor!")

	// after cleanup, signal to the room that the game state is done cleaning up
	g.cleanedUpChan <- true
}

func (g *GameStateProcessor) Cleanup() <-chan bool {
	return g.cleanedUpChan
}

func (g *GameStateProcessor) Status() model.GameStatus {
	return g.status.Status()
}

func (g *GameStateProcessor) StartGame() bool {
	// Check all players are ready
	playerOrderIDs, ok := g.players.AllPlayersReady()
	if !ok {
		log.Error().Msg("Attempted to start the game but not all players are ready!")
		return false
	}

	// Sanity check that the length of ready users is <= room's max capacity
	numPlayersReady := len(playerOrderIDs)
	if numPlayersReady > g.players.MaxPlayers() {
		log.Error().
			Int("numPlayersReady", numPlayersReady).
			Int("maxPlayers", g.players.MaxPlayers()).
			Msg("Number of ready players somehow exceeds the room's max capacity!")
		return false
	}

	// Randomize player turn order
	rand.Shuffle(numPlayersReady, func(i, j int) {
		playerOrderIDs[i], playerOrderIDs[j] = playerOrderIDs[j], playerOrderIDs[i]
	})

	// Save turn order
	g.status.SetPlayerOrderIDs(playerOrderIDs)

	// Get turn order
	currentPlayerTurn, err := g.getCurrentTurnUser()
	if err != nil {
		log.Error().
			Msg("Player order was computed but the current user turn references an invalid user ID")
		return false
	}

	// Notify all users that the game has started
	g.players.BroadcastEvent(events.StartGame(playerOrderIDs, currentPlayerTurn))

	// Set status to game started
	g.status.SetStatus(model.StatusStarted)

	// Progress game state logic with timer
	go g.gameLoop()

	return true
}

func (g *GameStateProcessor) HandleUserConnection(ctx context.Context, user *user.User, connErrChan chan error) {
	// Save user connection
	player := g.players.SaveConnection(user)

	// Concurrently handle the user's WebSocket connection
	go user.ReaderLoop(ctx, g.messageQueue, connErrChan)
	go user.WriterLoop(ctx, connErrChan)

	userModel := model.User{
		ID:   player.ID(),
		Name: player.Name(),
	}

	// Get current turn user model as a pointer
	var currentTurnUserPtr *model.User = nil
	currentTurnUser, err := g.getCurrentTurnUser()
	if err == nil {
		currentTurnUserPtr = &currentTurnUser
	}

	// Send rehydration event to user who just joined
	player.SendMessage(
		events.RehydrateForUser(
			userModel,
			g.players.PlayersAsModelList(),
			g.chatHistory.GetAll(),
			g.status.Status(),
			g.players.MaxPlayers(),
			g.status.PlayerOrderIDs(),
			currentTurnUserPtr,
			g.drawing.GetAll(),
		),
	)

	// Broadcast that a user has joined
	g.players.BroadcastEvent(events.UserJoin(player.ToModel(g.players.RoomLeaderID())))
	userJoinChatEvent := events.ChatSystemEvent(userModel.Name + " has joined the room.")
	g.sendChatAll(userJoinChatEvent)
}

func (g *GameStateProcessor) RemoveUserConnection(userID string) {
	// Remove user connection
	player := g.players.RemoveConnection(userID)

	// Broadcast that a user has left
	if player != nil {
		userModel := player.ToUserModel()
		g.players.BroadcastEvent(events.UserLeave(player.ToModel(g.players.RoomLeaderID())))
		userLeftChatEvent := events.ChatSystemEvent(userModel.Name + " has left the room.")
		g.sendChatAll(userLeftChatEvent)
	}
}
