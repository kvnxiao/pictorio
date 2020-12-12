package user

import (
	"context"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
)

type User struct {
	outgoing chan []byte
	conn     *websocket.Conn
	ID       string
	Name     string
}

func NewUser(conn *websocket.Conn, player model.User) *User {
	return &User{
		outgoing: make(chan []byte),
		conn:     conn,
		ID:       player.ID,
		Name:     player.Name,
	}
}

// ReaderLoop represents the read-loop that continuously ingests new messages from a user's WebSocket connection.
func (p *User) ReaderLoop(ctx context.Context, messageQueue chan<- []byte, connErrChan chan<- error) {
	for {
		_, readBytes, err := p.conn.Read(ctx)
		if err != nil {
			connErrChan <- err
			log.Info().Msg("DONE reader!")
			return
		}
		log.Info().
			Bytes("msg", readBytes).
			Str("id", p.ID).
			Msg("Received message from user")
		messageQueue <- readBytes
	}
}

// WriterLoop represents the write-loop that continuously ingests messages queued into the user's outgoing message
// channel and writes to the user's WebSocket connection.
func (p *User) WriterLoop(ctx context.Context, connErrChan chan error) {
	for {
		select {
		case msg := <-p.outgoing:
			err := p.conn.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				connErrChan <- err
				return
			}
			log.Info().
				Bytes("msg", msg).
				Str("uid", p.ID).
				Msg("Wrote message to user")
		case <-ctx.Done():
			log.Info().Msg("DONE writer!")
			return
		}
	}
}

// Outgoing returns the writable []byte channel that can be used to send messages to this specific player
func (p *User) Outgoing() chan<- []byte {
	return p.outgoing
}
