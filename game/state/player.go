package state

import (
	"github.com/kvnxiao/pictorio/game/user"
)

type PlayerState interface {
	Points() int
	Wins() int
	IsSpectator() bool
	IsConnected() bool
	IsReady() bool

	SetNewConnection(user *user.User)
	SetConnected(connected bool)
	SetReady(ready bool)
}

type Player struct {
	user        *user.User
	points      int
	wins        int
	isSpectator bool
	isConnected bool
	isReady     bool
}

func NewPlayer(user *user.User, isSpectator bool) PlayerState {
	return &Player{
		user:        user,
		points:      0,
		wins:        0,
		isSpectator: isSpectator,
		isConnected: false,
		isReady:     false,
	}
}

func (p *Player) Points() int {
	return p.points
}

func (p *Player) Wins() int {
	return p.wins
}

func (p *Player) IsConnected() bool {
	return p.isConnected
}

func (p *Player) IsSpectator() bool {
	return p.isSpectator
}

func (p *Player) IsReady() bool {
	return p.isReady
}

func (p *Player) SetNewConnection(user *user.User) {
	p.user = user
}

func (p *Player) SetConnected(connected bool) {
	p.isConnected = connected
}

func (p *Player) SetReady(ready bool) {
	p.isReady = ready
}
