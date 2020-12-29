package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type DrawTempEvent struct {
	User model.User `json:"user"`
	Line model.Line `json:"line"`
}

type DrawTempStopEvent struct {
	User model.User `json:"user"`
	Line model.Line `json:"line"`
}

func (e DrawTempEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e DrawTempEvent) GameEventType() GameEventType {
	return EventTypeDrawTemp
}

func (e DrawTempStopEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e DrawTempStopEvent) GameEventType() GameEventType {
	return EventTypeDrawTempStop
}

