package state

import (
	"math/rand"
	"strings"
	"time"

	"github.com/kvnxiao/pictorio/events"
	"github.com/kvnxiao/pictorio/game/settings"
	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

// gameLoop describes the game loop
// When game starts:
//   -> First drawer's turn gets some short duration to choose a randomly generated word
//   -> First drawer begins their drawing with a 60 second time limit
//      |-> Other users guess while first drawer is drawing
//   -> Award points if drawer guesses drawing correctly
func (g *GameStateProcessor) gameLoop() {
	// while game is in started state, continue game loop
	for g.status.Status() == model.StatusStarted {
		g.NextTurn()
	}
}

func (g *GameStateProcessor) waitForSelectedWord(
	currentTurnUser model.User,
	words []string,
	maxTimeSeconds int,
) string {
	timeLeftSeconds := maxTimeSeconds
	timeout := time.After(time.Duration(maxTimeSeconds) * time.Second)
	ticker := time.Tick(1 * time.Second)

	var selectedWord string

	startTime := time.Now().UnixNano()
	for {
		select {
		case <-timeout:
			// Player did not select a word in time, auto select a word for them
			selectedWord = words[rand.Intn(len(words))]
			return selectedWord
		case <-ticker:
			// Send player a decremented TurnCountdown event
			timeLeftSeconds -= 1
			if timeLeftSeconds >= 0 {
				g.players.BroadcastEvent(events.TurnCountdown(currentTurnUser, timeLeftSeconds))
			}
		case selectionIndex := <-g.wordSelectionIndex:
			// Ignore elements from selection index channel if the timestamp is before when startTime was calculated
			if selectionIndex.Timestamp >= startTime {
				if selectionIndex.Value >= len(words) {
					log.Error().
						Msg("Word selection index out of bounds, exceeds the number of generated random words")
				} else {
					selectedWord = words[selectionIndex.Value]
					return selectedWord
				}
			}
		}
	}
}

func (g *GameStateProcessor) handleGuess(
	currentTurnUser model.User,
	chosenWord string,
	censoredWord string,
	guessedPlayers map[string]bool,
	guess Guess,
) {
	if chosenWord == strings.TrimSpace(guess.Value) {
		if currentTurnUser.ID == guess.User.ID || guessedPlayers[guess.User.ID] {
			// Send censored word if user has already guessed the word, or the drawer is trying to send the word
			g.sendChatAll(events.ChatUserEvent(guess.User, censoredWord))
		} else {
			// First time the user is guessing the word correctly
			guessedPlayers[guess.User.ID] = true
			g.sendChatAll(events.ChatSystemEvent(guess.User.Name + " has guessed the word."))
			// TODO: award player with points
		}
	} else {
		g.sendChatAll(events.ChatUserEvent(guess.User, guess.Value))
	}
}

func (g *GameStateProcessor) waitForGuessOrTimeout(currentTurnUser model.User, maxTimeSeconds int) {
	guessedPlayers := make(map[string]bool)
	currentWord := g.status.CurrentWord()
	censoredWord := strings.Repeat("*", len(currentWord))

	timeLeftSeconds := maxTimeSeconds
	timeout := time.After(time.Duration(maxTimeSeconds) * time.Second)
	ticker := time.Tick(1 * time.Second)

	startTime := time.Now().UnixNano()
	for {
		select {
		case <-timeout:
			// fail-safe timeout has been reached
			return
		case <-ticker:
			timeLeftSeconds -= 1
			if timeLeftSeconds < 0 {
				// end turn if no more time remaining
				return
			}
		case guess := <-g.wordGuess:
			// Ignore elements from word guess channel if the timestamp is before when startTime was calculated
			if guess.Timestamp >= startTime {
				g.handleGuess(currentTurnUser, currentWord, censoredWord, guessedPlayers, guess)
			}
		}
	}
}

func (g *GameStateProcessor) generateWordLengths(word string) []int {
	split := strings.Fields(word)

	wordLengths := make([]int, len(split))
	for i := 0; i < len(split); i++ {
		wordLengths[i] = len(split[i])
	}
	return wordLengths
}

// beginTurnSelection starts the turn selection
func (g *GameStateProcessor) beginTurnSelection(userModel model.User) ([]string, int) {
	g.status.SetTurnStatus(model.TurnSelection)

	// Generate random word list (words that have not been recorded yet)
	words := g.status.GenerateWords()
	maxSelectionTimeSeconds := settings.MaxTurnSelectionCountdownSeconds

	// Send TurnBeginSelection event to current turn player (with the words)
	// Send TurnBeginSelection event to the other players (without the words)
	g.players.SendEvent(events.TurnBeginSelectionCurrentPlayer(userModel, maxSelectionTimeSeconds, words), userModel.ID)
	g.players.BroadcastEventExclude(events.TurnBeginSelection(userModel, maxSelectionTimeSeconds), userModel.ID)

	return words, maxSelectionTimeSeconds
}

// beginTurnDrawing starts the turn drawing
func (g *GameStateProcessor) beginTurnDrawing(userModel model.User, wordLengths []int, word string) int {
	g.status.SetTurnStatus(model.TurnDrawing)

	maxDrawingTimeSeconds := settings.MaxTurnDrawingCountdownSeconds

	// Send TurnBeginDrawing event to current turn player (with the selected word)
	// Send TurnBeginDrawing event to the other players (without the selected word)
	g.players.SendEvent(
		events.TurnBeginDrawingCurrentPlayer(userModel, maxDrawingTimeSeconds, wordLengths, word), userModel.ID,
	)
	g.players.BroadcastEventExclude(
		events.TurnBeginDrawing(userModel, maxDrawingTimeSeconds, wordLengths), userModel.ID,
	)

	return maxDrawingTimeSeconds
}

func (g *GameStateProcessor) endTurn(userModel model.User) {
	g.status.SetTurnStatus(model.TurnEnded)
	g.players.BroadcastEvent(events.TurnOver(userModel))
}

// NextTurn begins the next turn by allowing the current turn's user to select a word with a time limit
func (g *GameStateProcessor) NextTurn() {
	log.Info().Msg("Starting next turn!")
	defer log.Info().Msg("Turn ended")

	// Get current turn user
	userModel, err := g.getCurrentTurnUser()
	if err != nil {
		log.Error().Msg("Failed to get current turn's player state")
		return
	}

	// Begin turn selection
	words, maxSelectionTimeSeconds := g.beginTurnSelection(userModel)

	// Wait for word selection from current turn player, and save the current word
	word := g.waitForSelectedWord(userModel, words, maxSelectionTimeSeconds)
	wordLengths := g.generateWordLengths(word)
	g.status.SetCurrentWord(word)

	// Begin turn drawing
	maxDrawingTimeSeconds := g.beginTurnDrawing(userModel, wordLengths, word)

	// Wait for player guesses, or timeout from current turn player drawing
	g.waitForGuessOrTimeout(userModel, maxDrawingTimeSeconds)

	// End current turn
	g.endTurn(userModel)
}
