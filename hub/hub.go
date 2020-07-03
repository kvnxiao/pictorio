package hub

import (
	"github.com/kvnxiao/pictorio/game"
)

// Hub keeps track a mapping of roomID strings to their respective rooms
type Hub struct {
	rooms map[string]*game.Room
}

// New constructs a new Hub with defaults.
func New() *Hub {
	return &Hub{
		rooms: make(map[string]*game.Room),
	}
}

func (h *Hub) generateUniqueID() string {
	for {
		roomID := game.GenerateRoomID()
		_, exists := h.rooms[roomID]
		if !exists {
			return roomID
		}
	}
}

func (h *Hub) NewRoom() *game.Room {
	roomID := h.generateUniqueID()
	r := game.NewRoom(roomID)
	h.rooms[roomID] = r
	return r
}

func (h *Hub) Room(roomID string) (*game.Room, bool) {
	r, ok := h.rooms[roomID]
	return r, ok
}
