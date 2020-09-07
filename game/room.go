package game

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/kvnxiao/pictorio/game/player"
	"github.com/kvnxiao/pictorio/random"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
)

type Room struct {
	// mu is a mutex for checking the state of the room, i.e. whether it is closed or not when a person joins
	mu sync.Mutex

	// closed tells when to stop accepting new WebSocket connections, to prevent new people from joining the room
	closed bool

	// cleanup is a channel that stops the eventLoop's messageQueue goroutine from running
	cleanup chan bool

	// roomID represents the unique identifier associated with the room
	roomID string

	// players represents a set of individual player pointers
	players *player.Players

	// messageQueue is the message queue of incoming messages from players
	// TODO: change channel type from []byte to a struct that contains the player information as well
	messageQueue chan []byte
}

// NewRoom creates an empty room with the provided roomID string and sets up the global.
func NewRoom(roomID string) *Room {
	ro := &Room{
		roomID:       roomID,
		players:      player.NewContainer(),
		messageQueue: make(chan []byte),
		cleanup:      make(chan bool),
		closed:       false,
	}
	go ro.eventLoop()
	return ro
}

// ID returns the unique room ID representing this room.
func (r *Room) ID() string {
	return r.roomID
}

// addPlayer registers a player who has joined the room.
func (r *Room) addPlayer(p *player.Player) {
	log.Info().Msg("Adding player")
	r.players.Add(p)
}

// removePlayer removes a player from the room.
func (r *Room) removePlayer(p *player.Player) {
	log.Info().Str("id", p.ID).Msg("Removing player")
	r.players.Remove(p)
}

// Count returns the number of players currently in the room.
func (r *Room) Count() int {
	return r.players.Count()
}

// Cleanup sends a clean-up signal to the running eventLoop which stops handling new messages, and also sets the room's
// closed state to true, so that the room will not accept new WebSocket connections.
func (r *Room) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cleanup <- true
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

	conn, err := websocket.Accept(w, req, nil)
	if err != nil {
		log.Err(err).Send()
		return
	}

	err = r.newPlayer(req.Context(), conn)
	if errors.Is(err, context.Canceled) {
		log.Err(err).Str("type", "cancelled").Send()
		return
	}
	if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		log.Info().Err(err).Str("type", "closed").Send()
		return
	}
}

// newPlayer instantiates a new player struct from the incoming WebSocket connection.
// This method blocks after a player is added to the room, and waits until an error is encountered from either reading
// from the player's WebSocket connection, or when the server fails to write to the player's connection.
func (r *Room) newPlayer(ctx context.Context, conn *websocket.Conn) error {
	errChan := make(chan error)

	// TODO: pass in real user id
	p := player.New(conn, random.RandString(12))

	r.addPlayer(p)
	defer r.removePlayer(p)

	log.Info().Str("id", p.ID).Msg("Added new player")
	go p.ReaderLoop(ctx, r.messageQueue, errChan)
	go p.WriterLoop(ctx, errChan)

	// blocks on waiting for an error to be sent to the errChan.
	// an error will be sent through the errChan if a player's connection fails to be read from,
	// or fails to be written to
	err := <-errChan
	return err
}

// eventLoop represents a single instance of (i.e. the current room's) game logic, which handles and
// processes incoming WebSocket messages from the player, as well as handles cleaning up the room when all players
// have left the room.
func (r *Room) eventLoop() {
	for {
		select {
		case msg := <-r.messageQueue:
			r.players.Broadcast(msg)
		case <-r.cleanup:
			log.Info().Msg("CLEANING UP ROOM!")
			return
		}
	}
}
