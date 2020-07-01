package gameserver

import (
	"github.com/kvnxiao/pictorio/room"
)

// GameServer keeps track a mapping of roomID strings to their respective rooms
type GameServer struct {
	rooms map[string]*room.Room
}

// New constructs a GameServer with defaults.
func New() *GameServer {
	return &GameServer{
		rooms: make(map[string]*room.Room),
	}
}

func (gs *GameServer) NewRoom() *room.Room {
	roomID := room.GenerateID()
	r := room.NewRoom(roomID)
	go r.Handle()
	gs.rooms[roomID] = r
	return r
}

func (gs *GameServer) Room(roomID string) (*room.Room, bool) {
	r, ok := gs.rooms[roomID]
	return r, ok
}
