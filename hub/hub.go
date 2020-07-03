package hub

import (
	"sync"
	"time"

	"github.com/kvnxiao/pictorio/game"
	"github.com/rs/zerolog/log"
)

// Hub keeps track a mapping of roomID strings to their respective rooms
type Hub struct {
	roomMu sync.Mutex
	rooms  map[string]*game.Room
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
	h.roomMu.Lock()
	h.rooms[roomID] = r
	h.roomMu.Unlock()

	go h.roomCleanupListener(r)

	return r
}

func (h *Hub) roomCleanupListener(r *game.Room) {
	everyMinute := time.NewTicker(1 * time.Minute)
	for range everyMinute.C {
		log.Info().Msg("test")
		if r.Count() == 0 {
			r.Cleanup()
			break
		}
	}
	log.Info().Str("roomID", r.ID()).Msg("Removing empty room from hub.")
	h.Remove(r.ID())
}

func (h *Hub) Room(roomID string) (*game.Room, bool) {
	h.roomMu.Lock()
	r, ok := h.rooms[roomID]
	h.roomMu.Unlock()
	return r, ok
}

func (h *Hub) Remove(roomID string) {
	h.roomMu.Lock()
	delete(h.rooms, roomID)
	h.roomMu.Unlock()
}
