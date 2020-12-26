package guess

import (
	"github.com/kvnxiao/pictorio/model"
)

// PlayerGuesses represents the state of guesses for a single drawing turn (the drawer, and the rest of the players).
// The first player to guess correctly is awarded the most points, and the drawer is awarded some points (slightly less
// than the first guesser), while the rest of the players will be awarded less points if guessing correctly after the
// first guess has been established
type PlayerGuesses struct {
	guessesRemaining map[string]struct{}
	maxGuesses       int
}

func NewPlayerGuesses(currentTurnUser model.User, connectedPlayers []model.User) *PlayerGuesses {
	// Create a set of players who have NOT guessed the word correctly yet
	playersNotGuessed := make(map[string]struct{})
	for _, player := range connectedPlayers {
		playersNotGuessed[player.ID] = struct{}{}
	}
	// Remove the drawer from the set since they do not participate in guessing the word
	delete(playersNotGuessed, currentTurnUser.ID)

	return &PlayerGuesses{
		guessesRemaining: playersNotGuessed,
		maxGuesses:       len(playersNotGuessed),
	}
}

func (g *PlayerGuesses) FinishedGuessing() bool {
	return len(g.guessesRemaining) == 0
}

func (g *PlayerGuesses) HasGuessed(playerID string) bool {
	_, hasNotGuessed := g.guessesRemaining[playerID]
	return !hasNotGuessed
}

func (g *PlayerGuesses) AddGuessed(playerID string) (guesserPoints int, drawerPoints int) {
	delete(g.guessesRemaining, playerID)
	if len(g.guessesRemaining) == g.maxGuesses-1 {
		return 3, 2
	}
	return 1, 0
}
