package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

// TurnEndEvent is the server-sourced event that notifies all players that the current turn has ended and a new turn
// will begin
//
// Only sent after the current turn player has already begun drawing (never sent during word selection turn status)
type TurnEndEvent struct {
	User model.User `json:"user"`
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
