package state

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/kvnxiao/pictorio/events"
	"github.com/kvnxiao/pictorio/model"
	"github.com/kvnxiao/pictorio/words"
	"github.com/rs/zerolog/log"
)

// gameLoop describes the game loop
//
// Game loop summary:
//   1. Get current turn's player state (the player who is drawing)
//   -> If current drawer is in the game as a player but is disconnected during this process, skip this player's turn
//      ^go to 6.
//   2. Notify which player is next up as the drawer
//      a. Send TurnNextPlayer event nonce: next player, max time
//      b. Send countdown for this turn state
//   3. Begin word selection
//      -> Drawer has some short duration to choose a randomly generated word
//      a. Send TurnWordSelection event nonce: current drawer, list of words, max time
//      b. Send countdown for this turn state
//   4. Wait for word selection
//      Either:
//        -> Drawer selects a word from the list of randomly generated words
//        -> Drawer waits for timeout
//   5. Begin drawing for the current drawer
//      -> Drawer is given a time limit to draw
//      a. Send TurnDrawing event nonce: current drawer, word lengths, word for the drawer, max time
//      b. Send countdown for this turn state
//   6. Wait for guesses or drawer to timeout
//      -> Hides the chosen word from chat with censored text (asterisks, e.g. '***')
//   	-> Award points if other players guesses drawing correctly
//      -> Censors the chat message for players who have already guessed the word correctly
//      -> Also censors the chat message from the drawer if they try to "cheat" and type out their chosen word
//   7. Notify end of current turn
//      -> Sets the next drawer's turn
//      -> Increments the round counter if the next turn loops back around to the first player
//      a. Send TurnEnd event nonce: the drawer, the word answer, max time
//      b. Send countdown for this turn state
//   8. End game loop if round counter reaches max rounds
func (g *GameStateProcessor) gameLoop() {
	// while game is in started state, continue game loop
	for g.status.Status() == model.GameStarted {
		log.Info().Msg("Starting next turn!")

		// 1. Send the next drawer, skip if disconnected during this check
		userModel, isConnected, err := g.getDrawerPlayer()
		if err != nil {
			log.Error().Msg("Failed to get current turn's player state")
			return
		}
		if !isConnected {
			if userModel.ID == "" {
				log.Error().Msg("Current drawer is not connected, but user ID was invalid")
				return
			}
			log.Info().Msg("Current drawer not connected! Skipping turn.")
			g.skipTurn()
			continue
		}

		// 2. Begin next player turn notification
		g.beginTurnNextPlayer(userModel)
		log.Info().Msg(userModel.Name + "'s turn starts.")

		// 3. Begin word selection
		generatedWords, maxSelectionTimeSeconds := g.beginWordSelection(userModel)

		// 4. Wait for word selection
		selectedWord := g.waitForSelectedWord(generatedWords, maxSelectionTimeSeconds)
		word := words.NewGameWord(selectedWord)
		g.status.SetCurrentWord(word)

		// 5. Begin turn drawing
		maxDrawingTimeSeconds := g.beginTurnDrawing(userModel, word)

		// 6. Wait for player guesses, or timeout from current drawer drawing
		g.waitForGuessOrTimeout(userModel, maxDrawingTimeSeconds)

		// 7. End current turn
		g.beginTurnEnd(userModel)

		log.Info().Msg("Turn ended")

		// 8. Check rounds to end game loop
		g.checkRounds()
	}
	g.gameOver()
}

func (g *GameStateProcessor) getDrawerPlayer() (model.User, bool, error) {
	currentTurnID := g.status.CurrentTurnID()
	player, ok := g.players.GetPlayer(currentTurnID)
	if !ok {
		return model.User{}, false, errors.New("could not get player state with invalid user id for current turn")
	}

	return player.ToUserModel(), player.IsConnected(), nil
}

func (g *GameStateProcessor) beginTurnNextPlayer(userModel model.User) {
	g.status.SetTurnStatus(model.TurnNextPlayer)

	maxTimeSeconds := g.status.MaxNextUpTimeSeconds()
	g.broadcast(events.TurnBeginNextPlayer(userModel, maxTimeSeconds))

	timeLeftSeconds := maxTimeSeconds
	timeout := time.After(time.Duration(maxTimeSeconds) * time.Second)
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			timeLeftSeconds -= 1
			g.status.SetTimeRemaining(timeLeftSeconds)
			g.broadcast(events.TurnNextPlayerCountdown(maxTimeSeconds, timeLeftSeconds))
		case <-timeout:
			g.status.SetTimeRemaining(0)
			g.broadcast(events.TurnNextPlayerCountdown(maxTimeSeconds, 0))
			return
		}
	}
}

func (g *GameStateProcessor) waitForSelectedWord(
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
			g.status.SetTimeRemaining(0)
			log.Debug().Msg("Timeout in selecting a word, randomly choosing word")
			selectedWord = words[rand.Intn(len(words))]
			return selectedWord
		case <-ticker:
			// Send player a decremented TurnCountdown event
			timeLeftSeconds -= 1
			g.status.SetTimeRemaining(timeLeftSeconds)
			if timeLeftSeconds >= 0 {
				log.Debug().Int("timeLeft", timeLeftSeconds).Msg("Counting down for selection")
				g.broadcast(events.TurnWordSelectionCountdown(maxTimeSeconds, timeLeftSeconds))
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
			guesser, ok := g.players.GetPlayer(guess.User.ID)
			if !ok {
				log.Error().Msg("Player guessed the word correctly but does not exist in the players list")
				return
			}
			drawer, _ := g.players.GetPlayer(currentTurnUser.ID)
			guessedPlayers[guess.User.ID] = true
			g.awardPoints(guesser, 100, drawer, 20)
			g.broadcastChat(events.ChatSystemEvent(guess.User.Name + " has guessed the word."))
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
			log.Debug().Msg("Fail-safe timeout in drawing")
			g.status.SetTimeRemaining(0)
			g.broadcast(events.TurnDrawingCountdown(maxTimeSeconds, 0))
			return
		case <-ticker:
			timeLeftSeconds -= 1
			g.status.SetTimeRemaining(timeLeftSeconds)
			if timeLeftSeconds <= 0 {
				g.broadcast(events.TurnDrawingCountdown(maxTimeSeconds, 0))
				return
			} else {
				log.Debug().Int("timeLeft", timeLeftSeconds).Msg("Counting down for drawing")
				g.broadcast(events.TurnDrawingCountdown(maxTimeSeconds, timeLeftSeconds))
			}
		case guess := <-g.wordGuess:
			// Ignore elements from word guess channel if the timestamp is before when startTime was calculated
			if guess.Timestamp >= startTime {
				g.handleGuess(currentTurnUser, currentWord, guessedPlayers, guess)
			}
		}
	}
}

// beginWordSelection starts the word selection process for the current turn
func (g *GameStateProcessor) beginWordSelection(userModel model.User) ([]string, int) {
	g.status.SetTurnStatus(model.TurnSelection)

	// Generate random word list (words that have not been recorded yet)
	generatedWords := g.status.GenerateWords()
	maxSelectionTimeSeconds := g.status.MaxSelectionTimeSeconds()
	g.status.SetTimeRemaining(maxSelectionTimeSeconds)

	// Send TurnBeginSelection event to current drawer (with the words)
	// Send TurnBeginSelection event to the other players (without the words)
	g.emit(events.TurnBeginSelectionCurrentPlayer(userModel, maxSelectionTimeSeconds, generatedWords), userModel.ID)
	g.broadcastExcluding(events.TurnBeginSelection(userModel, maxSelectionTimeSeconds), userModel.ID)

	return generatedWords, maxSelectionTimeSeconds
}

// beginTurnDrawing starts the turn drawing
func (g *GameStateProcessor) beginTurnDrawing(userModel model.User, word words.GameWord) int {
	g.status.SetTurnStatus(model.TurnDrawing)

	maxDrawingTimeSeconds := g.status.MaxTurnDrawingTimeSeconds()
	g.status.SetTimeRemaining(maxDrawingTimeSeconds)

	// Send TurnBeginDrawing event to current drawer (with the selected word)
	// Send TurnBeginDrawing event to the other players (without the selected word)
	g.emit(
		events.TurnBeginDrawingCurrentPlayer(userModel, maxDrawingTimeSeconds, word.WordLength(), word.Word()),
		userModel.ID,
	)
	g.broadcastExcluding(events.TurnBeginDrawing(userModel, maxDrawingTimeSeconds, word.WordLength()), userModel.ID)

	return maxDrawingTimeSeconds
}

func (g *GameStateProcessor) skipTurn() {
	g.status.IncrementNextTurn()
}

func (g *GameStateProcessor) beginTurnEnd(userModel model.User) {
	g.status.SetTurnStatus(model.TurnEnded)
	word := g.status.CurrentWord().Word()

	// Notify that the current drawer's turn is ending, and broadcast what the word was
	maxTurnEndTimeSeconds := g.status.MaxTurnEndTimeSeconds()
	g.broadcast(events.TurnBeginEnd(userModel, word, maxTurnEndTimeSeconds))

	// Clear drawing state
	g.drawingHistory.Clear()

	// Increment current turn to the next user,
	// this will also will increment the round counter if the next turn loops back to first player
	g.status.IncrementNextTurn()

	timeLeftSeconds := maxTurnEndTimeSeconds
	timeout := time.After(time.Duration(maxTurnEndTimeSeconds) * time.Second)
	ticker := time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			timeLeftSeconds -= 1
			g.status.SetTimeRemaining(timeLeftSeconds)
			g.broadcast(events.TurnEndCountdown(maxTurnEndTimeSeconds, timeLeftSeconds))
		case <-timeout:
			g.status.SetTimeRemaining(0)
			g.broadcast(events.TurnEndCountdown(maxTurnEndTimeSeconds, 0))
			return
		}
	}
}

// checkRounds returns a boolean of whether the rounds played has exceeded the maximum number of rounds to be played
func (g *GameStateProcessor) checkRounds() bool {
	log.Info().Int("round", g.status.CurrentRound()).Msg("Current round is.")
	if g.status.CurrentRound() >= g.status.MaxRounds() {
		g.status.SetStatus(model.GameOver)
		return true
	}
	return false
}

func (g *GameStateProcessor) gameOver() {
	// TODO: game over
	log.Info().Msg("Game over!")
}
