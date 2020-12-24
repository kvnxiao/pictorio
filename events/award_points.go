package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type AwardPointsEvent struct {
	User   model.User `json:"user"`
	Points int        `json:"points"`
}

func (e AwardPointsEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e AwardPointsEvent) GameEventType() GameEventType {
	return EventTypeAwardPoints
}
