package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type DrawEventType int

const (
	Line DrawEventType = iota
	Undo
	Redo
	Clear
)

type DrawEvent struct {
	User model.User    `json:"user"`
	Line *model.Line   `json:"line,omitempty"`
	Type DrawEventType `json:"type"`
}

func (e DrawEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e DrawEvent) GameEventType() GameEventType {
	return EventTypeDraw
}
