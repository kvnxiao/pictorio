package state

import (
	"context"
	"encoding/json"
	"math/rand"

	"github.com/kvnxiao/pictorio/events"
	"github.com/kvnxiao/pictorio/game/settings"
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

	HandleUserConnection(ctx context.Context, user *user.User, connErrChan chan error)
	RemoveUserConnection(userID string)
}

// GameStateProcessor handles the state of the game
type GameStateProcessor struct {
	// roomID is the room ID associated with this game state
	roomID string

	// status is the current Status of the game
	status status.GameStatus

	// players represents the userID -> player states mapping
	players players.Players

	// drawingHistory is the current drawing history of the game
	drawingHistory drawing.History

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

func NewGameStateProcessor(roomID string) GameState {
	s := settings.DefaultSettings()

	return &GameStateProcessor{
		roomID:             roomID,
		status:             status.NewGameStatus(s),
		players:            players.NewPlayerContainer(s.MaxPlayers),
		drawingHistory:     drawing.NewDrawingHistory(),
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
			case events.EventTypeRehydrate:
			case events.EventTypeStartGame:
			case events.EventTypeTurnWordSelection:
			case events.EventTypeTurnDrawing:
			case events.EventTypeTurnEnd:
				g.warnServerSourcedEvent(event.Type)

			case events.EventTypeChat:
				var chatEvent events.ChatEvent
				err := json.Unmarshal(event.Data, &chatEvent)
				if err != nil {
					log.Error().Err(err).
						Msg("Could not unmarshal " + event.Type.String() + " from user")
				}
				g.onChatEvent(chatEvent)

			case events.EventTypeDraw:
				var drawEvent events.DrawEvent
				err := json.Unmarshal(event.Data, &drawEvent)
				if err != nil {
					log.Error().Err(err).
						Msg("Could not unmarshal " + event.Type.String() + " from user")
				}
				g.onDrawEvent(drawEvent)

			case events.EventTypeReady:
				var readyEvent events.ReadyEvent
				err := json.Unmarshal(event.Data, &readyEvent)
				if err != nil {
					log.Error().Err(err).
						Msg("Could not unmarshal " + event.Type.String() + " from user")
				}
				g.onReadyEvent(readyEvent)

			case events.EventTypeStartGameIssued:
				var startGameIssuedEvent events.StartGameIssuedEvent
				err := json.Unmarshal(event.Data, &startGameIssuedEvent)
				if err != nil {
					log.Error().Err(err).
						Msg("Could not unmarshal " + event.Type.String() + " from user")
				}
				g.onStartGameIssuedEvent(startGameIssuedEvent)

			case events.EventTypeTurnWordSelected:
				var turnWordSelectedEvent events.TurnWordSelectedEvent
				err := json.Unmarshal(event.Data, &turnWordSelectedEvent)
				if err != nil {
					log.Error().Err(err).
						Msg("Could not unmarshal " + event.Type.String() + " from user")
				}
				g.onTurnWordSelectedEvent(turnWordSelectedEvent)

			case events.EventTypeNewGameIssued:
				var newGameIssuedEvent events.NewGameIssuedEvent
				err := json.Unmarshal(event.Data, &newGameIssuedEvent)
				if err != nil {
					log.Error().Err(err).
						Msg("Could not unmarshal " + event.Type.String() + " from user")
				}
				g.onNewGameIssued(newGameIssuedEvent)

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
	log.Info().Str("roomID", g.roomID).Msg("Cleaning up game state processor")
	defer log.Info().Str("roomID", g.roomID).Msg("Done cleaning up game state processor!")

	// Cleanup game state processor
	g.chatHistory.Clear()
	g.drawingHistory.Clear()
	g.status.Reset()
	g.players.Cleanup()
	g.chatHistory = nil
	g.drawingHistory = nil
	g.status = nil
	g.players = nil
	close(g.messageQueue)
	close(g.wordSelectionIndex)
	close(g.wordGuess)

	// TODO: close channels if necessary

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

	// Notify all users that the game has started
	g.broadcast(events.StartGameEvent{PlayerOrderIDs: playerOrderIDs})

	// Set status to game started
	g.status.SetStatus(model.GameStarted)

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
	var selfUserIsCurrentTurn = false
	currentTurnUser, _, err := g.getDrawerPlayer()
	if err == nil {
		// A current turn exists, meaning the game has already started (otherwise it would be null)
		selfUserIsCurrentTurn = currentTurnUser.ID == userModel.ID
		currentTurnUserPtr = &currentTurnUser
	}

	// Send rehydration event to user who just joined
	g.emit(
		events.RehydrateForUser(
			userModel,
			currentTurnUserPtr,
			g.chatHistory.GetAll(),
			g.players.Summary(),
			g.status.Summary(selfUserIsCurrentTurn),
			g.drawingHistory.GetAll(),
		),
		userModel.ID,
	)

	// Broadcast user joined
	g.broadcast(events.UserJoin(player.ToModel(g.players.RoomLeaderID())))
	g.broadcastChat(events.ChatUserJoined(userModel))
}

func (g *GameStateProcessor) RemoveUserConnection(userID string) {
	// Remove user connection
	player := g.players.RemoveConnection(userID)

	// Broadcast user left
	if player != nil {
		userModel := player.ToUserModel()
		g.broadcast(events.UserLeave(player.ToModel(g.players.RoomLeaderID())))
		g.broadcastChat(events.ChatUserLeft(userModel))
	}
}

func (g *GameStateProcessor) awardPoints(
	guesser players.PlayerState, guesserPoints int,
	drawer players.PlayerState, drawerPoints int,
) {
	guesser.AwardPoints(guesserPoints)
	drawer.AwardPoints(drawerPoints)
	g.broadcast(events.AwardPointsEvent{
		Guesser:       guesser.ToUserModel(),
		Drawer:        drawer.ToUserModel(),
		GuesserPoints: guesserPoints,
		DrawerPoints:  drawerPoints,
	})
}
