package player

import (
	"context"

	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
)

type Player struct {
	outgoing chan []byte
	conn     *websocket.Conn
	ID       string
}

func New(conn *websocket.Conn, id string) *Player {
	return &Player{
		outgoing: make(chan []byte),
		conn:     conn,
		ID:       id,
	}
}

// ReaderLoop represents the read-loop that continuously ingests new messages from a player's WebSocket connection.
func (p *Player) ReaderLoop(ctx context.Context, messageQueue chan<- []byte, errChan chan<- error) {
	for {
		_, readBytes, err := p.conn.Read(ctx)
		if err != nil {
			errChan <- err
			log.Info().Msg("DONE reader!")
			return
		} else {
			log.Info().Bytes("msg", readBytes).Str("id", p.ID).Msg("Read something from player")
		}
		messageQueue <- readBytes
	}
}

// WriterLoop represents the write-loop that continuously ingests messages queued into the player's outgoing message
// channel and writes to the player's WebSocket connection.
func (p *Player) WriterLoop(ctx context.Context, errChan chan error) {
	for {
		select {
		case msg := <-p.outgoing:
			err := p.conn.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				errChan <- err
				return
			}
		case <-ctx.Done():
			log.Info().Msg("DONE writer!")
			return
		}
	}
}
