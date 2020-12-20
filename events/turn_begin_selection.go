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

func (e TurnBeginSelectionEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e TurnBeginSelectionEvent) GameEventType() GameEventType {
	return EventTypeTurnBeginSelection
}

func turnBeginSelection(currentTurnUser model.User, maxTimeSeconds int, words []string) TurnBeginSelectionEvent {
	return TurnBeginSelectionEvent{
		User:    currentTurnUser,
		MaxTime: maxTimeSeconds,
		Words:   words,
	}
}

func TurnBeginSelection(currentTurnUser model.User, maxTimeSeconds int) TurnBeginSelectionEvent {
	return turnBeginSelection(currentTurnUser, maxTimeSeconds, nil)
}

func TurnBeginSelectionCurrentPlayer(
	currentTurnUser model.User,
	maxTimeSeconds int,
	words []string,
) TurnBeginSelectionEvent {
	return turnBeginSelection(currentTurnUser, maxTimeSeconds, words)
}
