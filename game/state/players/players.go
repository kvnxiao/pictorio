package players

import (
	"sync"

	"github.com/kvnxiao/pictorio/game/user"
	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type Players interface {
	MaxPlayers() int
	RoomLeaderID() string

	GetPlayer(userID string) (PlayerState, bool)
	PlayersAsModelList() []model.PlayerState

	ReadyPlayer(userID string, ready bool) bool
	AllPlayersReady() ([]string, bool)

	SaveConnection(u *user.User) PlayerState
	RemoveConnection(userID string) PlayerState

	BroadcastEvent(eventBytes []byte)
	BroadcastEventExclude(eventBytes []byte, userID string)
	SendEvent(eventBytes []byte, userID string)
}

type PlayerStatesMap struct {
	mu           sync.RWMutex
	maxPlayers   int
	players      map[string]PlayerState
	roomLeaderID string
}

func NewPlayerContainer(maxPlayers int) Players {
	return &PlayerStatesMap{
		maxPlayers:   maxPlayers,
		players:      make(map[string]PlayerState),
		roomLeaderID: "",
	}
}

func (s *PlayerStatesMap) MaxPlayers() int {
	return s.maxPlayers
}

func (s *PlayerStatesMap) RoomLeaderID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.roomLeaderID
}

func (s *PlayerStatesMap) GetPlayer(userID string) (PlayerState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	player, ok := s.players[userID]
	return player, ok
}

func (s *PlayerStatesMap) PlayersAsModelList() []model.PlayerState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var playerStates []model.PlayerState
	for _, p := range s.players {
		playerStates = append(playerStates, p.ToModel(s.roomLeaderID))
	}

	return playerStates
}

func (s *PlayerStatesMap) ReadyPlayer(userID string, ready bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	player, ok := s.players[userID]
	if !ok {
		log.Error().Msg("Attempted to change ready state on an invalid user ID")
	}

	player.SetReady(ready)
	return ready
}

func (s *PlayerStatesMap) AllPlayersReady() ([]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var playerOrderIDs []string
	for _, player := range s.players {
		if player.IsConnected() && !player.IsSpectator() {
			if !player.IsReady() {
				return nil, false
			}
			playerOrderIDs = append(playerOrderIDs, player.ID())
		}
	}
	return playerOrderIDs, true
}

func (s *PlayerStatesMap) SaveConnection(u *user.User) PlayerState {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.roomLeaderID == "" {
		s.roomLeaderID = u.ID
	}

	player, ok := s.players[u.ID]
	if ok {
		// Existing user has re-joined the room
		player.SetNewConnection(u)
	} else {
		// New player has joined the room
		isRoomFull := len(s.players) >= s.maxPlayers
		player = newPlayer(u, isRoomFull)
		s.players[u.ID] = player
	}
	player.SetConnected(true)
	return player
}

func (s *PlayerStatesMap) RemoveConnection(userID string) PlayerState {
	s.mu.Lock()
	defer s.mu.Unlock()

	player, ok := s.players[userID]
	if !ok {
		return nil
	}

	player.SetConnected(false)
	player.SetReady(false)

	return player
}

func (s *PlayerStatesMap) BroadcastEvent(eventBytes []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, player := range s.players {
		if player.IsConnected() {
			player.SendMessage(eventBytes)
		}
	}
}

func (s *PlayerStatesMap) BroadcastEventExclude(eventBytes []byte, userID string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, player := range s.players {
		if player.IsConnected() && player.ID() != userID {
			player.SendMessage(eventBytes)
		}
	}
}

func (s *PlayerStatesMap) SendEvent(eventBytes []byte, userID string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	player, ok := s.players[userID]
	if !ok {
		log.Error().Msg("Attempted to send an event to an invalid player ID")
	}
	player.SendMessage(eventBytes)
}
