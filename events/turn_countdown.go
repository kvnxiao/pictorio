package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

// TurnCountdownEvent is the server-sourced event that simply counts down the number of seconds the current turn player
// has left to complete their turn action
type TurnCountdownEvent struct {
	User     model.User `json:"user"`
	TimeLeft int        `json:"timeLeft"`
}

func (e TurnCountdownEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e TurnCountdownEvent) GameEventType() GameEventType {
	return EventTypeTurnCountdown
}
