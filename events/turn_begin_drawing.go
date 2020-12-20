package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

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

func (e TurnBeginDrawingEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e TurnBeginDrawingEvent) GameEventType() GameEventType {
	return EventTypeTurnBeginDrawing
}

func turnBeginDrawing(
	currentTurnUser model.User,
	maxTimeSeconds int,
	wordLengths []int,
	word string,
) TurnBeginDrawingEvent {
	var wordPtr *string = nil
	if word != "" {
		wordPtr = &word
	}

	return TurnBeginDrawingEvent{
		User:       currentTurnUser,
		MaxTime:    maxTimeSeconds,
		WordLength: wordLengths,
		Word:       wordPtr,
	}
}

func TurnBeginDrawing(currentTurnUser model.User, maxTimeSeconds int, wordLengths []int) TurnBeginDrawingEvent {
	return turnBeginDrawing(currentTurnUser, maxTimeSeconds, wordLengths, "")
}

func TurnBeginDrawingCurrentPlayer(
	currentTurnUser model.User,
	maxTimeSeconds int,
	wordLengths []int,
	word string,
) TurnBeginDrawingEvent {
	return turnBeginDrawing(currentTurnUser, maxTimeSeconds, wordLengths, word)
}
