package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type TurnDrawingNonce struct {
	User       model.User `json:"user"`
	WordLength []int      `json:"wordLength"`
	Word       *string    `json:"word,omitempty"`
}

// TurnDrawingEvent is the server-sourced event that notifies all players the current turn player has selected a
// word and is beginning their drawing
//
// To current turn player:
// - Word is non-nil
//
// To rest of players:
// - Word is nil
type TurnDrawingEvent struct {
	Nonce    *TurnDrawingNonce `json:"nonce,omitempty"`
	Hints    []model.Hint      `json:"hints,omitempty"`
	MaxTime  int               `json:"maxTime"`
	TimeLeft int               `json:"timeLeft"`
	Status   model.TurnStatus  `json:"status"`
}

func (e TurnDrawingEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e TurnDrawingEvent) GameEventType() GameEventType {
	return EventTypeTurnDrawing
}

func turnBeginDrawing(
	currentTurnUser model.User,
	maxTimeSeconds int,
	wordLengths []int,
	word string,
) TurnDrawingEvent {
	var wordPtr *string = nil
	if word != "" {
		wordPtr = &word
	}

	return TurnDrawingEvent{
		Nonce: &TurnDrawingNonce{
			User:       currentTurnUser,
			WordLength: wordLengths,
			Word:       wordPtr,
		},
		Hints:    nil,
		MaxTime:  maxTimeSeconds,
		TimeLeft: maxTimeSeconds,
		Status:   model.TurnDrawing,
	}
}

func TurnBeginDrawing(currentTurnUser model.User, maxTimeSeconds int, wordLengths []int) TurnDrawingEvent {
	return turnBeginDrawing(currentTurnUser, maxTimeSeconds, wordLengths, "")
}

func TurnBeginDrawingCurrentPlayer(
	currentTurnUser model.User,
	maxTimeSeconds int,
	wordLengths []int,
	word string,
) TurnDrawingEvent {
	return turnBeginDrawing(currentTurnUser, maxTimeSeconds, wordLengths, word)
}

func TurnDrawingCountdown(maxTimeSeconds, timeLeftSeconds int, hints []model.Hint) TurnDrawingEvent {
	return TurnDrawingEvent{
		Nonce:    nil,
		Hints:    hints,
		MaxTime:  maxTimeSeconds,
		TimeLeft: timeLeftSeconds,
		Status:   model.TurnDrawing,
	}
}
