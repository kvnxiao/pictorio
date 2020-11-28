package player

import (
	"sync"

	"github.com/rs/zerolog/log"
)

type Players struct {
	// mu is a mutex for synchronizing on reads and modifications to the players map
	mu sync.Mutex

	// players represents a map of individual player ids to a Player pointer
	players map[string]*Player
}

func NewContainer() *Players {
	return &Players{
		players: make(map[string]*Player),
	}
}

func (l *Players) Add(p *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.players[p.ID.String()] = p
}

func (l *Players) Remove(p *Player) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.players, p.ID.String())
}

func (l *Players) Count() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return len(l.players)
}

func (l *Players) Broadcast(msg []byte) {
	for _, p := range l.players {
		p.outgoing <- msg
	}
}

func (l *Players) Send(playerID string, msg []byte) {
	p, ok := l.players[playerID]
	if ok {
		log.Info().Msg(string(msg))
		p.outgoing <- msg
	}
}
