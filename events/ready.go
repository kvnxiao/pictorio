package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type ReadyEvent struct {
	User  model.User `json:"user"`
	Ready bool       `json:"ready"`
}

func (e ReadyEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e ReadyEvent) GameEventType() GameEventType {
	return EventTypeReady
}
