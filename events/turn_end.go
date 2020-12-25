package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type TurnEndNonce struct {
	User   model.User `json:"user"`
	Answer string     `json:"answer"`
}

// TurnEndEvent is the server-sourced event that notifies all players that the current turn has ended and a new turn
// will begin
//
// Only sent after the current turn player has already begun drawing (never sent during word selection turn status)
type TurnEndEvent struct {
	Nonce    *TurnEndNonce    `json:"nonce,omitempty"`
	MaxTime  int              `json:"maxTime"`
	TimeLeft int              `json:"timeLeft"`
	Status   model.TurnStatus `json:"status"`
}

func (e TurnEndEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e TurnEndEvent) GameEventType() GameEventType {
	return EventTypeTurnEnd
}

func TurnBeginEnd(userModel model.User, word string, maxTimeSeconds int) TurnEndEvent {
	return TurnEndEvent{
		Nonce: &TurnEndNonce{
			User:   userModel,
			Answer: word,
		},
		MaxTime:  maxTimeSeconds,
		TimeLeft: maxTimeSeconds,
		Status:   model.TurnEnded,
	}
}

func TurnEndCountdown(maxTimeSeconds int, timeLeftSeconds int) TurnEndEvent {
	return TurnEndEvent{
		Nonce:    nil,
		MaxTime:  maxTimeSeconds,
		TimeLeft: timeLeftSeconds,
		Status:   model.TurnEnded,
	}
}
