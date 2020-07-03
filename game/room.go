package game

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
)

type player struct {
	outgoing chan []byte
	conn     *websocket.Conn
	id       string
}

type Room struct {
	roomID   string
	playerMu sync.Mutex
	players  map[*player]struct{}

	// mq is the message queue of incoming messages from players
	// TODO: change channel type from []byte to a struct that contains the player information as well
	mq chan []byte

	// cleanup is a channel that stops the globalWriter message queue goroutine from running
	cleanup chan bool
	closed  bool
}

func NewRoom(roomID string) *Room {
	ro := &Room{
		roomID:  roomID,
		players: make(map[*player]struct{}),
		mq:      make(chan []byte),
		cleanup: make(chan bool),
	}
	go ro.globalWriter()
	return ro
}

func (r *Room) ID() string {
	return r.roomID
}

// addPlayer registers a player who has joined the room
func (r *Room) addPlayer(p *player) {
	r.playerMu.Lock()
	r.players[p] = struct{}{}
	r.playerMu.Unlock()
}

// removePlayer removes a player from the room
func (r *Room) removePlayer(p *player) {
	log.Info().Msg("Removing player")
	r.playerMu.Lock()
	delete(r.players, p)
	log.Info().Int("count", len(r.players)).Send()
	r.playerMu.Unlock()
}

// count returns the number of players currently in the room
func (r *Room) Count() int {
	r.playerMu.Lock()
	defer r.playerMu.Unlock()
	return len(r.players)
}

// ConnectionHandler accepts a new WebSocket connection from the http request, and then subscribes it to all
// future messages
func (r *Room) ConnectionHandler(w http.ResponseWriter, req *http.Request) {
	if r.closed {
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

func (r *Room) newPlayer(ctx context.Context, conn *websocket.Conn) error {
	errChan := make(chan error)

	p := &player{
		outgoing: make(chan []byte),
		conn:     conn,
	}

	r.addPlayer(p)
	defer r.removePlayer(p)

	log.Info().Msg("Added new player")
	go r.reader(ctx, p, errChan)
	go r.writer(ctx, p, errChan)

	err := <-errChan
	return err
}

func (r *Room) reader(ctx context.Context, p *player, errChan chan error) {
	for {
		_, b, err := p.conn.Read(ctx)
		log.Info().Err(err).Bytes("msg", b).Msg("Read something new!")
		if err != nil {
			errChan <- err
			log.Error().Msg("DONE reader!")
			return
		}
		r.mq <- b
		log.Info().Bytes("msg", b)
	}
}

func (r *Room) writer(ctx context.Context, p *player, errChan chan error) {
	for {
		select {
		case msg := <-p.outgoing:
			err := p.conn.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				errChan <- err
				return
			}
		case <-ctx.Done():
			log.Error().Msg("DONE writer!")
			return
		}
	}
}

func (r *Room) globalWriter() {
	for {
		select {
		case msg := <-r.mq:
			for p := range r.players {
				if p != nil {
					p.outgoing <- msg
				}
			}
		case <-r.cleanup:
			log.Info().Msg("CLEANING UP ROOM!")
			return
		}
	}
}

func (r *Room) Cleanup() {
	r.cleanup <- true
	r.closed = true
}
