package events

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type GameEventType int

const (
	EventTypeUserJoinLeaveEvent GameEventType = iota
	EventTypeRehydrate
	EventTypeChat
	EventTypeDraw
)

type GameEvent struct {
	Type GameEventType   `json:"type"`
	Data json.RawMessage `json:"data"`
}

func gameEvent(eventType GameEventType, rawEventData json.RawMessage) []byte {
	if rawEventData == nil {
		return nil
	}
	gameEvent := GameEvent{
		Type: eventType,
		Data: rawEventData,
	}
	bytes, err := json.Marshal(&gameEvent)
	if err != nil {
		log.Err(err).Msg("Could not marshal GameEvent wrapper into JSON.")
		return nil
	}
	return bytes
}
