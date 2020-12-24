package status

import (
	"sync"

	"github.com/kvnxiao/pictorio/game/settings"
	"github.com/kvnxiao/pictorio/model"
	"github.com/kvnxiao/pictorio/words"
)

type GameStatus interface {
	Summary(selfUserIsCurrentTurn bool) model.GameStateSummary

	MaxRounds() int
	MaxSelectionTimeSeconds() int
	MaxTurnTimeSeconds() int
	CurrentRound() int

	Status() model.GameStatus
	SetStatus(status model.GameStatus)

	TurnStatus() model.TurnStatus
	SetTurnStatus(turnStatus model.TurnStatus)

	CurrentWord() words.GameWord
	SetCurrentWord(word words.GameWord)

	CurrentTurnID() string
	IncrementNextTurn() int

	PlayerOrderIDs() []string
	SetPlayerOrderIDs(playerOrderIDs []string)

	GenerateWords() []string
	WordSelections() []string

	SetTimeRemaining(seconds int)

	Cleanup()
}

type Status struct {
	mu             sync.RWMutex
	maxRounds      int
	currentRound   int
	status         model.GameStatus
	turnStatus     model.TurnStatus
	currentWord    words.GameWord
	playerOrderIDs []string
	turnIndex      int
	wordHistory    map[string]bool

	// Temporary
	maxSelectionTime int
	maxTurnTime      int
	timeLeftSeconds  int
	wordSelections   []string
}

func NewGameStatus(maxRounds int, maxSelectionSeconds int, maxTurnSeconds int) GameStatus {
	return &Status{
		// required fields
		maxRounds:        maxRounds,
		maxSelectionTime: maxSelectionSeconds,
		maxTurnTime:      maxTurnSeconds,
		currentRound:     1,
		status:           model.GameWaitingReadyUp,
		turnStatus:       model.TurnSelection,
		currentWord:      words.GameWord{},
		playerOrderIDs:   nil,
		turnIndex:        0,
		wordHistory:      make(map[string]bool),

		// initialize temp storage variables
		timeLeftSeconds: 0,
		wordSelections:  nil,
	}
}

func (s *Status) Summary(selfUserIsCurrentTurn bool) model.GameStateSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	wordSelections := s.wordSelections
	if s.turnStatus == model.TurnSelection {
		wordSelections = nil
	}

	word := s.currentWord.Word()
	if !selfUserIsCurrentTurn {
		word = ""
	}

	return model.GameStateSummary{
		MaxRounds:        s.maxRounds,
		MaxSelectionTime: s.maxSelectionTime,
		MaxTurnTime:      s.maxTurnTime,
		Round:            s.currentRound,
		TimeLeft:         s.timeLeftSeconds,
		Status:           s.status,
		TurnStatus:       s.turnStatus,
		PlayerOrderIDs:   s.playerOrderIDs,
		WordSummary: model.WordSummary{
			Word:           word,
			WordLength:     s.currentWord.WordLength(),
			WordSelections: wordSelections,
		},
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

func (s *Status) CurrentWord() words.GameWord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.currentWord
}

func (s *Status) SetCurrentWord(word words.GameWord) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.currentWord = word
	s.wordHistory[word.Word()] = true
}

func (s *Status) MaxRounds() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.maxRounds
}

func (s *Status) MaxSelectionTimeSeconds() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.maxSelectionTime
}

func (s *Status) MaxTurnTimeSeconds() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.maxTurnTime
}

func (s *Status) CurrentRound() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.currentRound
}

func (s *Status) CurrentTurnID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.playerOrderIDs) > 0 {
		return s.playerOrderIDs[s.turnIndex]
	}
	return ""
}

func (s *Status) IncrementNextTurn() int {
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

	var w []string
	for len(w) < settings.MaxSelectableWords {
		word := words.GenerateWord()
		if !s.wordHistory[word] {
			w = append(w, word)
		}
	}

	s.wordSelections = w

	return w
}

func (s *Status) WordSelections() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.wordSelections
}

func (s *Status) SetTimeRemaining(seconds int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.timeLeftSeconds = seconds
}

func (s *Status) Cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status = model.GameWaitingReadyUp
	s.turnStatus = model.TurnSelection
	s.currentWord = words.GameWord{}
	s.currentRound = 1
	s.turnIndex = 0
	s.playerOrderIDs = nil
	s.wordHistory = nil
}
