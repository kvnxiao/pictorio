package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type DrawSelectThicknessEvent struct {
	User           model.User `json:"user"`
	ThicknessIndex int        `json:"thicknessIdx"`
}

func (e DrawSelectThicknessEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e DrawSelectThicknessEvent) GameEventType() GameEventType {
	return EventTypeDrawSelectThickness
}
