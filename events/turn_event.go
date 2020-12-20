package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

// TurnBeginSelectionEvent is the server-sourced event in which the current turn player receives a list of randomly
// generated words that they can select to begin their drawing turn, and the rest of the players are sent the same
// event without any words
//
// To current turn player:
// - Words is a non-nil list of strings
//
// To rest of players:
// - Words is nil
type TurnBeginSelectionEvent struct {
	User    model.User `json:"user"`
	MaxTime int        `json:"maxTime"`
	Words   []string   `json:"words,omitempty"`
}

// TurnWordSelectedEvent is the client-sourced event in which the current turn player selects a word from the generated
// list of words that the server produces (this list is sent to the client from beforehand in a TurnEvent)
type TurnWordSelectedEvent struct {
	User  model.User `json:"user"`
	Index int        `json:"index"`
}

// TurnBeginDrawingEvent is the server-sourced event that notifies all players the current turn player has selected a
// word and is beginning their drawing
//
// To current turn player:
// - Word is non-nil
//
// To rest of players:
// - Word is nil
type TurnBeginDrawingEvent struct {
	User       model.User `json:"user"`
	MaxTime    int        `json:"maxTime"`
	WordLength []int      `json:"wordLength"`
	Word       *string    `json:"word,omitempty"`
}

// TurnCountdownEvent is the server-sourced event that simply counts down the number of seconds the current turn player
// has left to complete their turn action
type TurnCountdownEvent struct {
	User     model.User `json:"user"`
	TimeLeft int        `json:"timeLeft"`
}

// TurnOverEvent is the server-sourced event that notifies all players that the current turn has ended and a new turn
// will begin
//
// Only sent after the current turn player has already begun drawing (never sent during word selection turn status)
type TurnOverEvent struct {
	User model.User `json:"user"`
}

func turnBeginSelection(currentTurnUser model.User, maxTimeSeconds int, words []string) []byte {
	event := TurnBeginSelectionEvent{
		User:    currentTurnUser,
		MaxTime: maxTimeSeconds,
		Words:   words,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal " + EventTypeTurnBeginSelection.String() + " into JSON.")
		return nil
	}

	return gameEvent(EventTypeTurnBeginSelection, eventBytes)
}

func turnBeginDrawing(currentTurnUser model.User, maxTimeSeconds int, wordLengths []int, word string) []byte {
	var wordPtr *string = nil
	if word != "" {
		wordPtr = &word
	}

	event := TurnBeginDrawingEvent{
		User:       currentTurnUser,
		MaxTime:    maxTimeSeconds,
		WordLength: wordLengths,
		Word:       wordPtr,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal " + EventTypeTurnBeginDrawing.String() + " into JSON.")
		return nil
	}

	return gameEvent(EventTypeTurnBeginDrawing, eventBytes)
}

func TurnBeginSelection(currentTurnUser model.User, maxTimeSeconds int) []byte {
	return turnBeginSelection(currentTurnUser, maxTimeSeconds, nil)
}

func TurnBeginSelectionCurrentPlayer(currentTurnUser model.User, maxTimeSeconds int, words []string) []byte {
	return turnBeginSelection(currentTurnUser, maxTimeSeconds, words)
}

func TurnCountdown(currentTurnUser model.User, timeLeftSeconds int) []byte {
	event := TurnCountdownEvent{
		User:     currentTurnUser,
		TimeLeft: timeLeftSeconds,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal " + EventTypeTurnCountdown.String() + " into JSON.")
		return nil
	}

	return gameEvent(EventTypeTurnCountdown, eventBytes)
}

func TurnBeginDrawing(currentTurnUser model.User, maxTimeSeconds int, wordLengths []int) []byte {
	return turnBeginDrawing(currentTurnUser, maxTimeSeconds, wordLengths, "")
}

func TurnBeginDrawingCurrentPlayer(currentTurnUser model.User, maxTimeSeconds int, wordLengths []int, word string) []byte {
	return turnBeginDrawing(currentTurnUser, maxTimeSeconds, wordLengths, word)
}

func TurnOver(currentTurnUser model.User) []byte {
	event := TurnOverEvent{
		User: currentTurnUser,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal " + EventTypeTurnOver.String() + " into JSON.")
		return nil
	}

	return gameEvent(EventTypeTurnOver, eventBytes)
}
