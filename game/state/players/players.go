package players

import (
	"sort"
	"sync"

	"github.com/kvnxiao/pictorio/events"
	"github.com/kvnxiao/pictorio/game/user"
	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type Players interface {
	Summary() model.PlayersSummary

	MaxPlayers() int
	RoomLeaderID() string

	GetPlayer(userID string) (PlayerState, bool)
	GetConnectedPlayers(includeSpectator bool) []model.User

	ReadyPlayer(userID string, ready bool) bool
	AllPlayersReady() ([]string, bool)
	AllPlayersDisconnected() bool

	SaveConnection(u *user.User) PlayerState
	RemoveConnection(userID string) PlayerState

	SendEventToAll(event events.SerializableEvent)
	SendEventToAllExcept(event events.SerializableEvent, userID string)
	SendEventToUser(event events.SerializableEvent, userID string)

	Winners() []model.Winner

	Reset()
	Cleanup()
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

func (s *PlayerStatesMap) Summary() model.PlayersSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var playerStates []model.PlayerState
	for _, p := range s.players {
		playerStates = append(playerStates, p.ToModel(s.roomLeaderID))
	}

	return model.PlayersSummary{
		PlayerStates: playerStates,
		MaxPlayers:   s.maxPlayers,
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

func (s *PlayerStatesMap) GetConnectedPlayers(includeSpectator bool) []model.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var connectedUsers []model.User

	for _, player := range s.players {
		if player.IsConnected() &&
			(player.IsSpectator() && includeSpectator || !player.IsSpectator()) {
			connectedUsers = append(connectedUsers, player.ToUserModel())
		}
	}

	return connectedUsers
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

func (s *PlayerStatesMap) AllPlayersDisconnected() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, player := range s.players {
		if !player.IsSpectator() && player.IsConnected() {
			return false
		}
	}
	return true
}

func (s *PlayerStatesMap) SaveConnection(u *user.User) PlayerState {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Set the first user as the room leader
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

func (s *PlayerStatesMap) SendEventToAll(event events.SerializableEvent) {
	eventBytes := events.ToJson(event)

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, player := range s.players {
		if player.IsConnected() {
			player.SendMessage(eventBytes)
		}
	}
}

func (s *PlayerStatesMap) SendEventToAllExcept(event events.SerializableEvent, userID string) {
	eventBytes := events.ToJson(event)

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, player := range s.players {
		if player.IsConnected() && player.ID() != userID {
			player.SendMessage(eventBytes)
		}
	}
}

func (s *PlayerStatesMap) SendEventToUser(event events.SerializableEvent, userID string) {
	eventBytes := events.ToJson(event)

	s.mu.RLock()
	defer s.mu.RUnlock()

	player, ok := s.players[userID]
	if !ok {
		log.Error().Msg("Attempted to send an event to an invalid player ID")
	}
	if player.IsConnected() {
		player.SendMessage(eventBytes)
	}
}

func (s *PlayerStatesMap) Winners() []model.Winner {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var winners []model.Winner

	for _, player := range s.players {
		if !player.IsSpectator() {
			winners = append(winners, model.Winner{
				User:   player.ToUserModel(),
				Points: player.Points(),
			})
		}
	}

	sort.Slice(winners, func(i, j int) bool {
		return winners[i].Points > winners[j].Points
	})

	return winners
}

func (s *PlayerStatesMap) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, player := range s.players {
		player.SetReady(false)
		player.ResetPoints()
	}
}

func (s *PlayerStatesMap) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k := range s.players {
		delete(s.players, k)
	}
}
