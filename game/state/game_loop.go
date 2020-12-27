package state

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/kvnxiao/pictorio/events"
	"github.com/kvnxiao/pictorio/game/guess"
	"github.com/kvnxiao/pictorio/game/hint"
	"github.com/kvnxiao/pictorio/game/settings"
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
		// Ensure that everyone currently playing is still connected before continuing
		if g.players.AllPlayersDisconnected() {
			log.Info().Msg("Seems like no one is left in the room for a game still in progress. Cleaning up.")
			g.status.SetStatus(model.GameOver)
			return
		}

		log.Debug().Msg("Starting next turn!")
		setting := g.status.Settings()

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
			log.Debug().Msg("Current drawer not connected! Skipping turn.")
			g.skipTurn()
			continue
		}

		// 2. Begin next player turn notification
		g.beginTurnNextPlayer(userModel, setting)

		// 3. Begin word selection
		// 4. Wait for word selection
		generatedWords, maxSelectionTimeSeconds := g.beginWordSelection(userModel, setting)
		selectedWord := g.waitForSelectedWord(generatedWords, maxSelectionTimeSeconds)
		word := words.NewGameWord(selectedWord)
		g.status.SetCurrentWord(word)

		// 5. Begin turn drawing
		// 6. Wait for player guesses, or timeout from current drawer drawing
		maxDrawingTimeSeconds := g.beginTurnDrawing(userModel, word, setting)
		g.waitForGuessOrTimeout(userModel, maxDrawingTimeSeconds, setting)

		// 7. End current turn
		g.beginTurnEnd(userModel, setting)

		// 8. Check rounds to end game loop
		if g.checkRounds(setting) {
			break
		}
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

func (g *GameStateProcessor) beginTurnNextPlayer(userModel model.User, setting settings.GameSettings) {
	log.Debug().Str("uid", userModel.ID).Msg("Beginning next turn phase for next player")
	g.status.SetTurnStatus(model.TurnNextPlayer)

	maxTimeSeconds := setting.MaxTurnNextPlayerTimeSeconds
	g.broadcast(events.TurnBeginNextPlayer(userModel, maxTimeSeconds))

	timeLeftSeconds := maxTimeSeconds
	timeout := time.After(time.Duration(maxTimeSeconds+1) * time.Second)
	ticker := time.Tick(1 * time.Second)
	for {
		select {

		case <-timeout:
			log.Debug().Str("uid", userModel.ID).Msg("Next turn is starting")
			g.status.SetTimeRemaining(0)
			g.broadcast(events.TurnNextPlayerCountdown(maxTimeSeconds, 0))
			return

		case <-ticker:
			timeLeftSeconds -= 1
			if timeLeftSeconds < 0 {
				timeLeftSeconds = 0
			}
			log.Debug().
				Str("uid", userModel.ID).
				Int("timeLeft", timeLeftSeconds).
				Msg("Next turn timer countdown")
			g.status.SetTimeRemaining(timeLeftSeconds)
			g.broadcast(events.TurnNextPlayerCountdown(maxTimeSeconds, timeLeftSeconds))
		}
	}
}

func (g *GameStateProcessor) waitForSelectedWord(
	words []string,
	maxTimeSeconds int,
) string {
	log.Debug().Msg("Waiting for word selection from drawer")

	timeLeftSeconds := maxTimeSeconds
	timeout := time.After(time.Duration(maxTimeSeconds+1) * time.Second)
	ticker := time.Tick(1 * time.Second)

	var selectedWord string

	startTime := time.Now().UnixNano()
	for {
		select {

		// Player did not select a word in time, auto select a word for them
		case <-timeout:
			log.Debug().
				Msg("Word selection timeout")

			g.status.SetTimeRemaining(0)
			selectedWord = words[rand.Intn(len(words))]
			return selectedWord

		// Send players a decrementing TurnWordSelection event
		case <-ticker:
			timeLeftSeconds -= 1
			if timeLeftSeconds < 0 {
				timeLeftSeconds = 0
			}

			log.Debug().
				Int("timeLeft", timeLeftSeconds).
				Msg("Word selection timer countdown")

			g.status.SetTimeRemaining(timeLeftSeconds)
			g.broadcast(events.TurnWordSelectionCountdown(maxTimeSeconds, timeLeftSeconds))

		case selectionIndex := <-g.wordSelectionIndex:
			// Ignore elements from selection index channel if the timestamp is before when startTime was calculated
			if selectionIndex.Timestamp >= startTime {
				if selectionIndex.Value >= len(words) {
					log.Error().
						Msg("Word selection index out of bounds, exceeds the number of generated random words")
				} else {
					selectedWord = words[selectionIndex.Value]
					log.Debug().
						Str("selectedWord", selectedWord).
						Msg("A word has been selected by the drawer")
					return selectedWord
				}
			}
		}
	}
}

func (g *GameStateProcessor) handleGuess(
	currentTurnUser model.User,
	word words.GameWord,
	guesses *guess.PlayerGuesses,
	wordGuess Guess,
) bool {
	candidate := strings.ToLower(strings.TrimSpace(wordGuess.Value))

	// Handle word match
	if word.Word() == candidate {
		if currentTurnUser.ID == wordGuess.User.ID || guesses.HasGuessed(wordGuess.User.ID) {
			// Send censored word if user has already guessed the word, or the drawer is trying to send the word
			g.broadcastChat(events.ChatUserEvent(wordGuess.User, word.Censored()))
			return false
		}

		// First time the user is guessing the word correctly
		guesser, ok := g.players.GetPlayer(wordGuess.User.ID)
		if !ok {
			log.Error().Msg("Player guessed the word correctly but does not exist in the players list")
			return false
		}
		drawer, _ := g.players.GetPlayer(currentTurnUser.ID)
		guesserPoints, drawerPoints := guesses.AddGuessed(wordGuess.User.ID)
		g.awardPoints(guesser, guesserPoints, drawer, drawerPoints)
		g.broadcastChat(events.ChatSystemEvent(wordGuess.User.Name + " has guessed the word."))
		return true
	}

	// Handle non-exact-match messages
	if strings.Contains(candidate, word.Word()) &&
		(currentTurnUser.ID == wordGuess.User.ID || guesses.HasGuessed(wordGuess.User.ID)) {
		// Censor text that contains the word as a substring
		g.broadcastChat(events.ChatUserEvent(wordGuess.User, words.Censor(len(wordGuess.Value))))
	} else {
		// Regular chat messages
		g.broadcastChat(events.ChatUserEvent(wordGuess.User, wordGuess.Value))
	}
	return false
}

func (g *GameStateProcessor) waitForGuessOrTimeout(
	currentTurnUser model.User,
	maxTimeSeconds int,
	setting settings.GameSettings,
) {
	log.Debug().Msg("Waiting for guess or timeout for the drawing phase")

	currentWord := g.status.CurrentWord()

	// Guess
	guesses := guess.NewPlayerGuesses(currentTurnUser, g.players.GetConnectedPlayers(false))
	firstGuess := false

	// Hints
	hints := hint.NewHint(currentWord.Hints(), setting.HintSettings)
	var hintsToSend = make([]model.Hint, 0)

	// Timer
	timeLeftSeconds := maxTimeSeconds
	timeout := time.After(time.Duration(maxTimeSeconds+1) * time.Second)
	ticker := time.Tick(1 * time.Second)

	startTime := time.Now().UnixNano()
	for {
		select {

		// Turn ended
		case <-timeout:
			log.Debug().Msg("Drawing phase timeout")

			g.status.SetTimeRemaining(0)
			g.broadcast(events.TurnDrawingCountdown(maxTimeSeconds, 0, hintsToSend))
			return

		// Send players a decrementing TurnDrawing event
		case <-ticker:
			timeLeftSeconds -= 1
			if timeLeftSeconds < 0 {
				timeLeftSeconds = 0
			}

			log.Debug().
				Int("timeLeft", timeLeftSeconds).
				Msg("Drawing phase timer countdown")

			g.status.SetTimeRemaining(timeLeftSeconds)

			if !firstGuess {
				nextHint, hasNextHint := hints.NextHint(timeLeftSeconds)
				if hasNextHint {
					log.Debug().
						Int("wordIndex", nextHint.WordIndex).
						Int("charIndex", nextHint.CharIndex).
						Str("char", string(nextHint.Char)).
						Msg("Generating next hint")
					hintsToSend = append(hintsToSend, nextHint)
				}
			}
			g.broadcast(events.TurnDrawingCountdown(maxTimeSeconds, timeLeftSeconds, hintsToSend))

		case wordGuess := <-g.wordGuess:
			// Ignore elements from word guess channel if the timestamp is before when startTime was calculated
			if wordGuess.Timestamp >= startTime {
				if g.handleGuess(currentTurnUser, currentWord, guesses, wordGuess) {
					if guesses.FinishedGuessing() {
						log.Debug().Msg("Everyone has guessed the word")

						timeLeftSeconds = 0
						g.status.SetTimeRemaining(0)
						g.broadcast(events.TurnDrawingCountdown(maxTimeSeconds, timeLeftSeconds, hintsToSend))
						return
					} else if timeLeftSeconds > setting.MaxTurnDrawingTimeCutSeconds {
						log.Debug().Msg("First guess of the word, reducing countdown timer")

						firstGuess = true
						timeLeftSeconds = setting.MaxTurnDrawingTimeCutSeconds
						g.status.SetTimeRemaining(timeLeftSeconds)
						g.broadcast(events.TurnDrawingCountdown(maxTimeSeconds, timeLeftSeconds, hintsToSend))
					}
				}
			}
		}
	}
}

// beginWordSelection starts the word selection process for the current turn
func (g *GameStateProcessor) beginWordSelection(userModel model.User, setting settings.GameSettings) ([]string, int) {
	log.Debug().Str("uid", userModel.ID).Msg("Beginning word selection phase for the drawer")
	g.status.SetTurnStatus(model.TurnSelection)

	// Generate random word list (words that have not been recorded yet)
	generatedWords := g.status.GenerateWords()
	maxSelectionTimeSeconds := setting.MaxTurnSelectionTimeSeconds
	g.status.SetTimeRemaining(maxSelectionTimeSeconds)

	log.Debug().
		Str("uid", userModel.ID).
		Strs("words", generatedWords).
		Msg("Generated random words for the drawer")

	// Send TurnBeginSelection event to current drawer (with the words)
	// Send TurnBeginSelection event to the other players (without the words)
	g.emit(events.TurnBeginSelectionCurrentPlayer(userModel, maxSelectionTimeSeconds, generatedWords), userModel.ID)
	g.broadcastExcluding(events.TurnBeginSelection(userModel, maxSelectionTimeSeconds), userModel.ID)

	return generatedWords, maxSelectionTimeSeconds
}

// beginTurnDrawing starts the turn drawing
func (g *GameStateProcessor) beginTurnDrawing(userModel model.User, word words.GameWord, setting settings.GameSettings) int {
	log.Debug().Str("uid", userModel.ID).Msg("Beginning drawing phase for the drawer")
	g.status.SetTurnStatus(model.TurnDrawing)

	maxDrawingTimeSeconds := setting.MaxTurnDrawingTimeSeconds
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

func (g *GameStateProcessor) beginTurnEnd(userModel model.User, setting settings.GameSettings) {
	log.Debug().Str("uid", userModel.ID).Msg("Beginning turn end phase for the drawer")

	g.status.SetTurnStatus(model.TurnEnded)
	word := g.status.CurrentWord().Word()

	// Notify that the current drawer's turn is ending, and broadcast what the word was
	maxTimeSeconds := setting.MaxTurnEndTimeSeconds
	g.broadcast(events.TurnBeginEnd(userModel, word, maxTimeSeconds))

	timeLeftSeconds := maxTimeSeconds
	timeout := time.After(time.Duration(maxTimeSeconds+1) * time.Second)
	ticker := time.Tick(1 * time.Second)
	for {
		select {

		case <-timeout:
			log.Debug().Msg("Turn end phase timeout")

			g.status.SetTimeRemaining(0)
			g.broadcast(events.TurnEndCountdown(maxTimeSeconds, 0))

			// Clear drawing state
			g.drawingHistory.Clear()

			// Increment current turn to the next user,
			// this will also will increment the round counter if the next turn loops back to first player
			g.status.IncrementNextTurn()

			return

		case <-ticker:
			timeLeftSeconds -= 1
			if timeLeftSeconds <= 0 {
				timeLeftSeconds = 0
			}

			log.Debug().Int("timeLeft", timeLeftSeconds).Msg("Turn end timer countdown")

			g.status.SetTimeRemaining(timeLeftSeconds)
			g.broadcast(events.TurnEndCountdown(maxTimeSeconds, timeLeftSeconds))
		}
	}
}

// checkRounds returns a boolean of whether the rounds played has exceeded the maximum number of rounds to be played
func (g *GameStateProcessor) checkRounds(setting settings.GameSettings) bool {
	log.Debug().Int("round", g.status.CurrentRound()).Msg("Checking current round")
	if g.status.CurrentRound() >= setting.MaxRounds {
		return true
	}
	return false
}

func (g *GameStateProcessor) gameOver() {
	log.Debug().Msg("Game over!")

	g.status.SetStatus(model.GameOver)
	g.broadcast(events.GameOverEvent{
		Winners: g.players.Winners(),
	})
}
