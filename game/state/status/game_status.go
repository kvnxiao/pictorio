package status

import (
	"sync"

	"github.com/kvnxiao/pictorio/game/settings"
	"github.com/kvnxiao/pictorio/model"
	"github.com/kvnxiao/pictorio/random"
)

type GameStatus interface {
	Status() model.GameStatus
	SetStatus(status model.GameStatus)
	TurnStatus() model.TurnStatus
	SetTurnStatus(turnStatus model.TurnStatus)

	CurrentWord() string
	SetCurrentWord(word string)
	CurrentRound() int

	CurrentTurnID() string
	TurnIndex() int
	NextTurnIndex() int

	PlayerOrderIDs() []string
	SetPlayerOrderIDs(playerOrderIDs []string)

	GenerateWords() []string

	Cleanup()
}

type Status struct {
	mu             sync.RWMutex
	status         model.GameStatus
	turnStatus     model.TurnStatus
	currentWord    string
	currentRound   int
	playerOrderIDs []string
	turnIndex      int
	wordHistory    map[string]bool
}

func NewGameStatus() GameStatus {
	return &Status{
		status:         model.StatusWaitingReadyUp,
		turnStatus:     model.TurnSelection,
		currentWord:    "",
		currentRound:   1,
		playerOrderIDs: nil,
		turnIndex:      0,
		wordHistory:    make(map[string]bool),
	}
}

func (s *Status) Status() model.GameStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.status
}

func (s *Status) SetStatus(status model.GameStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status = status
}

func (s *Status) TurnStatus() model.TurnStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.turnStatus
}

func (s *Status) SetTurnStatus(turnStatus model.TurnStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.turnStatus = turnStatus
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
	s.wordHistory[word] = true
}

func (s *Status) CurrentRound() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.currentRound
}

func (s *Status) CurrentTurnID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.playerOrderIDs[s.turnIndex]
}

func (s *Status) TurnIndex() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.turnIndex
}

func (s *Status) NextTurnIndex() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	nextTurnIndex := (s.turnIndex + 1) % len(s.playerOrderIDs)
	if nextTurnIndex == 0 {
		s.currentRound += 1
	}

	s.turnIndex = nextTurnIndex
	return nextTurnIndex
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
	s.turnIndex = 0
}

func (s *Status) GenerateWords() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var words []string
	for len(words) < settings.MaxSelectableWords {
		word := random.GenerateWord()
		if !s.wordHistory[word] {
			words = append(words, word)
		}
	}

	return words
}

func (s *Status) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status = model.StatusWaitingReadyUp
	s.currentWord = ""
	s.currentRound = 1
	s.turnIndex = 0
	s.playerOrderIDs = nil
	s.wordHistory = nil
}
