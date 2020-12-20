package status

import (
	"sync"

	"github.com/kvnxiao/pictorio/model"
)

type GameStatus interface {
	Status() model.GameStateStatus
	SetStatus(status model.GameStateStatus)

	CurrentWord() string
	SetCurrentWord(word string)

	CurrentTurnID() string

	PlayerOrderIDs() []string
	SetPlayerOrderIDs(playerOrderIDs []string)

	Cleanup()
}

type Status struct {
	mu             sync.RWMutex
	status         model.GameStateStatus
	currentWord    string
	currentTurnID  string
	playerOrderIDs []string
}

func NewGameStatus() GameStatus {
	return &Status{
		status:         model.StatusWaitingReadyUp,
		currentWord:    "",
		currentTurnID:  "",
		playerOrderIDs: nil,
	}
}

func (s *Status) Status() model.GameStateStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.status
}

func (s *Status) SetStatus(status model.GameStateStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status = status
}

func (s *Status) CurrentWord() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.currentWord
}

func (s *Status) SetCurrentWord(word string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentWord = word
}

func (s *Status) CurrentTurnID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.currentTurnID
}

func (s *Status) PlayerOrderIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.playerOrderIDs
}

func (s *Status) SetPlayerOrderIDs(playerOrderIDs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.playerOrderIDs = playerOrderIDs
	s.currentTurnID = playerOrderIDs[0]
}

func (s *Status) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status = model.StatusWaitingReadyUp
	s.currentWord = ""
	s.currentTurnID = ""
	s.playerOrderIDs = nil
}
