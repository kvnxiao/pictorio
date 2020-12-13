package players

import (
	"github.com/kvnxiao/pictorio/game/user"
	"github.com/kvnxiao/pictorio/model"
)

type PlayerState interface {
	ID() string
	Name() string

	Points() int
	Wins() int
	IsSpectator() bool
	IsConnected() bool
	IsReady() bool
	IsRoomLeader(roomLeaderUserID string) bool

	SetNewConnection(u *user.User)
	SetConnected(connected bool)
	SetReady(ready bool)

	SendMessage(bytes []byte)

	ToModel(roomLeaderUserID string) model.PlayerState
	ToUserModel() model.User
}

type Player struct {
	user        *user.User
	points      int
	wins        int
	isSpectator bool
	isConnected bool
	isReady     bool
}

func newPlayer(u *user.User, isSpectator bool) PlayerState {
	return &Player{
		user:        u,
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

func (p *Player) ID() string {
	return p.user.ID
}

func (p *Player) Name() string {
	return p.user.Name
}

func (p *Player) SetNewConnection(u *user.User) {
	p.user = u
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

func (p *Player) ToModel(roomLeaderUserID string) model.PlayerState {
	return model.PlayerState{
		User: model.User{
			ID:   p.ID(),
			Name: p.Name(),
		},
		Points:       p.Points(),
		Wins:         p.Wins(),
		IsSpectator:  p.IsSpectator(),
		IsConnected:  p.IsConnected(),
		IsReady:      p.IsReady(),
		IsRoomLeader: p.IsRoomLeader(roomLeaderUserID),
	}
}

func (p *Player) ToUserModel() model.User {
	return model.User{
		ID:   p.ID(),
		Name: p.Name(),
	}
}
