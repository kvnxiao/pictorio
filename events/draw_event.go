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

func Draw(event DrawEvent) []byte {
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal " + EventTypeDraw.String() + " into JSON.")
		return nil
	}

	return gameEvent(EventTypeDraw, eventBytes)
}
