package state

import (
	"github.com/kvnxiao/pictorio/game/user"
	"github.com/kvnxiao/pictorio/model"
)

type PlayerState interface {
	Points() int
	Wins() int
	IsSpectator() bool
	IsConnected() bool
	IsReady() bool
	IsRoomLeader(roomLeaderUserID string) bool
	UserModel() model.User

	SetNewConnection(user *user.User)
	SetConnected(connected bool)
	SetReady(ready bool)

	SendMessage(bytes []byte)
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

func (p *Player) IsRoomLeader(roomLeaderUserID string) bool {
	return p.user.ID == roomLeaderUserID
}

func (p *Player) UserModel() model.User {
	return model.User{
		ID:   p.user.ID,
		Name: p.user.Name,
	}
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

func (p *Player) SendMessage(bytes []byte) {
	p.user.Outgoing() <- bytes
}
