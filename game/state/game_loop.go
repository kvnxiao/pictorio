package state

import (
	"math/rand"
	"strings"
	"time"

	"github.com/kvnxiao/pictorio/events"
	"github.com/kvnxiao/pictorio/game/settings"
	"github.com/kvnxiao/pictorio/model"
	"github.com/kvnxiao/pictorio/words"
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
	for g.status.Status() == model.GameStarted {
		g.nextTurn()
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
			log.Info().Msg("Timeout in selecting a word, randomly choosing word")
			selectedWord = words[rand.Intn(len(words))]
			return selectedWord
		case <-ticker:
			// Send player a decremented TurnCountdown event
			timeLeftSeconds -= 1
			if timeLeftSeconds >= 0 {
				log.Info().Int("timeLeft", timeLeftSeconds).Msg("Counting down for selection")
				g.broadcast(events.TurnCountdownEvent{User: currentTurnUser, TimeLeft: timeLeftSeconds})
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
	word words.GameWord,
	guessedPlayers map[string]bool,
	guess Guess,
) {
	candidate := strings.ToLower(strings.TrimSpace(guess.Value))

	if word.Word() == candidate {
		if currentTurnUser.ID == guess.User.ID || guessedPlayers[guess.User.ID] {
			// Send censored word if user has already guessed the word, or the drawer is trying to send the word
			g.broadcastChat(events.ChatUserEvent(guess.User, word.Censored()))
		} else {
			// First time the user is guessing the word correctly
			guessedPlayers[guess.User.ID] = true
			g.broadcastChat(events.ChatSystemEvent(guess.User.Name + " has guessed the word."))
			// TODO: award player with points
		}
	} else if strings.Contains(candidate, word.Word()) &&
		(currentTurnUser.ID == guess.User.ID || guessedPlayers[guess.User.ID]) {
		g.broadcastChat(events.ChatUserEvent(guess.User, words.Censor(len(guess.Value))))
	} else {
		g.broadcastChat(events.ChatUserEvent(guess.User, guess.Value))
	}
}

func (g *GameStateProcessor) waitForGuessOrTimeout(currentTurnUser model.User, maxTimeSeconds int) {
	guessedPlayers := make(map[string]bool)
	currentWord := g.status.CurrentWord()

	timeLeftSeconds := maxTimeSeconds
	timeout := time.After(time.Duration(maxTimeSeconds) * time.Second)
	ticker := time.Tick(1 * time.Second)

	startTime := time.Now().UnixNano()
	for {
		select {
		case <-timeout:
			// fail-safe timeout has been reached
			log.Info().Msg("Fail-safe timeout in drawing")
			return
		case <-ticker:
			timeLeftSeconds -= 1
			if timeLeftSeconds < 0 {
				// end turn if no more time remaining
				return
			} else {
				log.Info().Int("timeLeft", timeLeftSeconds).Msg("Counting down for drawing")
				g.broadcast(events.TurnCountdownEvent{User: currentTurnUser, TimeLeft: timeLeftSeconds})
			}
		case guess := <-g.wordGuess:
			// Ignore elements from word guess channel if the timestamp is before when startTime was calculated
			if guess.Timestamp >= startTime {
				g.handleGuess(currentTurnUser, currentWord, guessedPlayers, guess)
			}
		}
	}
}

// beginTurnSelection starts the turn selection
func (g *GameStateProcessor) beginTurnSelection(userModel model.User) ([]string, int) {
	g.status.SetTurnStatus(model.TurnSelection)

	// Generate random word list (words that have not been recorded yet)
	generatedWords := g.status.GenerateWords()
	maxSelectionTimeSeconds := settings.MaxTurnSelectionCountdownSeconds

	// Send TurnBeginSelection event to current turn player (with the words)
	// Send TurnBeginSelection event to the other players (without the words)
	g.emit(events.TurnBeginSelectionCurrentPlayer(userModel, maxSelectionTimeSeconds, generatedWords), userModel.ID)
	g.broadcastExcluding(events.TurnBeginSelection(userModel, maxSelectionTimeSeconds), userModel.ID)

	return generatedWords, maxSelectionTimeSeconds
}

// beginTurnDrawing starts the turn drawing
func (g *GameStateProcessor) beginTurnDrawing(userModel model.User, word words.GameWord) int {
	g.status.SetTurnStatus(model.TurnDrawing)

	maxDrawingTimeSeconds := settings.MaxTurnDrawingCountdownSeconds

	// Send TurnBeginDrawing event to current turn player (with the selected word)
	// Send TurnBeginDrawing event to the other players (without the selected word)
	g.emit(
		events.TurnBeginDrawingCurrentPlayer(userModel, maxDrawingTimeSeconds, word.WordLength(), word.Word()),
		userModel.ID,
	)
	g.broadcastExcluding(events.TurnBeginDrawing(userModel, maxDrawingTimeSeconds, word.WordLength()), userModel.ID)

	return maxDrawingTimeSeconds
}

func (g *GameStateProcessor) endTurn(userModel model.User) {
	g.status.SetTurnStatus(model.TurnEnded)
	g.broadcast(events.TurnEndEvent{User: userModel})
}

// NextTurn begins the next turn by allowing the current turn's user to select a word with a time limit
func (g *GameStateProcessor) nextTurn() {
	log.Info().Msg("Starting next turn!")
	defer log.Info().Msg("Turn ended")

	// Get current turn user
	userModel, err := g.getCurrentTurnUser()
	if err != nil {
		log.Error().Msg("Failed to get current turn's player state")
		return
	}

	// Begin turn selection
	generatedWords, maxSelectionTimeSeconds := g.beginTurnSelection(userModel)
	log.Info().Strs("words", generatedWords).Msg("Generated random words")

	// Wait for word selection from current turn player, and save the current word
	selectedWord := g.waitForSelectedWord(userModel, generatedWords, maxSelectionTimeSeconds)
	word := words.NewGameWord(selectedWord)
	log.Info().Str("word", word.Word()).Ints("wordLength", word.WordLength()).Msg("Word selected")

	g.status.SetCurrentWord(word)

	// Begin turn drawing
	maxDrawingTimeSeconds := g.beginTurnDrawing(userModel, word)

	// Wait for player guesses, or timeout from current turn player drawing
	g.waitForGuessOrTimeout(userModel, maxDrawingTimeSeconds)

	// End current turn
	g.endTurn(userModel)
}
