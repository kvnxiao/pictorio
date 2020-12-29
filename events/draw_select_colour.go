package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type DrawSelectColourEvent struct {
	User        model.User `json:"user"`
	ColourIndex int        `json:"colourIdx"`
}

func (e DrawSelectColourEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e DrawSelectColourEvent) GameEventType() GameEventType {
	return EventTypeDrawSelectColour
}
