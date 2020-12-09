package game

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/kvnxiao/pictorio/cookies"
	"github.com/kvnxiao/pictorio/ctxs"
	"github.com/kvnxiao/pictorio/game/state"
	"github.com/kvnxiao/pictorio/game/user"
	"github.com/kvnxiao/pictorio/model"
	"github.com/kvnxiao/pictorio/random"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/ksuid"
	"nhooyr.io/websocket"
)

const (
	roomIDLength = 9

	maxPlayerNum = 8
)

type Room struct {
	// roomID represents the unique identifier associated with the room
	roomID string

	// mu is a mutex for checking the state of the room, i.e. whether it is closed or not when a person joins
	mu sync.Mutex
	// userMu is a mutex for handling websocket connections between
	userMu sync.Mutex

	// closed tells when to stop accepting new WebSocket connections, to prevent new people from joining the room
	closed bool

	// usersMap represents a set of individual users mapped by their user-ID
	usersMap map[string]*user.User

	gameProcessor state.GameState

	// startCleanupChan is a channel that signals the gameProcessor's event loop to stop running
	startCleanupChan chan bool
}

// NewRoom creates an empty room with the provided roomID string and sets up the global.
func NewRoom(roomID string) *Room {
	room := &Room{
		roomID:           roomID,
		closed:           false,
		usersMap:         make(map[string]*user.User),
		gameProcessor:    state.NewGameStateProcessor(maxPlayerNum),
		startCleanupChan: make(chan bool),
	}
	go room.gameProcessor.EventProcessor(room.startCleanupChan)
	return room
}

// ID returns the unique room ID representing this room.
func (r *Room) ID() string {
	return r.roomID
}

// addUser registers a user who has joined the room.
func (r *Room) addUser(user *user.User) {
	r.userMu.Lock()
	defer r.userMu.Unlock()

	r.usersMap[user.ID] = user
}

// removeUser removes a user from the room.
func (r *Room) removeUser(user *user.User) {
	r.userMu.Lock()
	defer r.userMu.Unlock()

	delete(r.usersMap, user.ID)
}

// Count returns the number of users currently in the room.
func (r *Room) Count() int {
	r.userMu.Lock()
	defer r.userMu.Unlock()

	return len(r.usersMap)
}

// Cleanup sends a clean-up signal to the running eventLoop which stops handling new messages, and also sets the room's
// closed state to true, so that the room will not accept new WebSocket connections.
func (r *Room) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// get cleanup channel from game processor
	gameProcessorCleanup := r.gameProcessor.Cleanup()

	// signal to game processor that room is ready to be cleaned up
	r.startCleanupChan <- true

	// wait for game processor to be done cleaning up
	<-gameProcessorCleanup

	// set status of room to closed
	r.closed = true
}

func (r *Room) isClosed() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.closed
}

// ConnectionHandler accepts a new WebSocket connection from the http request, and then subscribes it to all
// future messages.
func (r *Room) ConnectionHandler(w http.ResponseWriter, req *http.Request) {
	if r.isClosed() {
		log.Error().Msg("Room is closed.")
		return
	}

	// Read user's unique ID or generate one if not exist
	userID, err := cookies.GetUserID(w, req)
	if err != nil || userID == "" {
		randomID, err := ksuid.NewRandom()
		if err != nil {
			log.Err(err).Msg("Could not generate a unique ID for new user")
			return
		}
		userID = randomID.String()
		cookies.SetUserID(w, userID)
	}

	// Ensure user ID parsed from cookies is of expected type
	userKSUID, err := ksuid.Parse(userID)
	if err != nil {
		log.Err(err).Msg("Could not parse user ID")
	}

	// Read user name
	userName, err := cookies.GetUserName(w, req)
	if err != nil || userName == "" {
		log.Info().Msg("Generating random name for new user")
		userName = random.GenerateName()
		cookies.SetUserName(w, userName)
	}

	// Save user ID and name to connection context
	ctx := context.WithValue(req.Context(), ctxs.KeyUserID, userKSUID)
	ctx = context.WithValue(ctx, ctxs.KeyUserName, userName)

	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		log.Err(err).Msg("Could not upgrade connection from user to a WebSocket connection.")
		return
	}

	err = r.newUserConnection(ctx, conn)
	if errors.Is(err, context.Canceled) {
		log.Err(err).Str("type", "cancelled").Msg("User connection closed.")
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		log.Info().Err(err).Str("type", "closed").Msg("User connection closed.")
		return
	}
}

// newUserConnection instantiates a new user struct from the incoming WebSocket connection.
// This method blocks after a user is added to the room, and waits until an error is encountered from either reading
// from the user's WebSocket connection, or when the server fails to write to the user's connection.
func (r *Room) newUserConnection(ctx context.Context, conn *websocket.Conn) error {
	connErrChan := make(chan error)

	userKSUID, ok := ctxs.UserID(ctx)
	if !ok {
		return errors.New("could not get user ID from connection context")
	}
	userName, ok := ctxs.UserName(ctx)
	if !ok {
		return errors.New("could not get user name from connection context")
	}

	userModel := model.User{
		ID:   userKSUID.String(),
		Name: userName,
	}

	u := user.NewUser(conn, userModel)

	r.addUser(u)
	defer r.removeUser(u)

	log.Info().
		Str("roomID", r.roomID).
		Str("uid", u.ID).
		Str("uname", u.Name).
		Msg("Added new user to room")

	r.gameProcessor.UserJoined(ctx, u, connErrChan)

	// blocks on waiting for an error to be sent to the connErrChan.
	// an error is sent through the channel if a user's connection either fails to be read from or written to
	err := <-connErrChan

	r.gameProcessor.UserLeft(userModel.ID)
	// user is removed after this function exits due to the defer statement

	return err
}
