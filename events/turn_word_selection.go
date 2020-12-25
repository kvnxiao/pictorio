package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type TurnWordSelectionNonce struct {
	User  model.User `json:"user"`
	Words []string   `json:"words,omitempty"`
}

// TurnWordSelectionEvent is the server-sourced event in which the current turn player receives a list of randomly
// generated words that they can select to begin their drawing turn, and the rest of the players are sent the same
// event without any words
//
// To current turn player:
// - Words is a non-nil list of strings
//
// To rest of players:
// - Words is nil
type TurnWordSelectionEvent struct {
	Nonce    *TurnWordSelectionNonce `json:"nonce,omitempty"`
	MaxTime  int                     `json:"maxTime"`
	TimeLeft int                     `json:"timeLeft"`
	Status   model.TurnStatus        `json:"status"`
}

func (e TurnWordSelectionEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e TurnWordSelectionEvent) GameEventType() GameEventType {
	return EventTypeTurnWordSelection
}

func turnBeginSelection(currentTurnUser model.User, maxTimeSeconds int, words []string) TurnWordSelectionEvent {
	return TurnWordSelectionEvent{
		Nonce: &TurnWordSelectionNonce{
			User:  currentTurnUser,
			Words: words,
		},
		MaxTime:  maxTimeSeconds,
		TimeLeft: maxTimeSeconds,
		Status:   model.TurnSelection,
	}
}

func TurnBeginSelection(currentTurnUser model.User, maxTimeSeconds int) TurnWordSelectionEvent {
	return turnBeginSelection(currentTurnUser, maxTimeSeconds, nil)
}

func TurnBeginSelectionCurrentPlayer(
	currentTurnUser model.User,
	maxTimeSeconds int,
	words []string,
) TurnWordSelectionEvent {
	return turnBeginSelection(currentTurnUser, maxTimeSeconds, words)
}

func TurnWordSelectionCountdown(maxTimeSeconds, timeLeftSeconds int) TurnWordSelectionEvent {
	return TurnWordSelectionEvent{
		Nonce:    nil,
		MaxTime:  maxTimeSeconds,
		TimeLeft: timeLeftSeconds,
		Status:   model.TurnSelection,
	}
}
