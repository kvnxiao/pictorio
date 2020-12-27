package events

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type SerializableEvent interface {
	RawJSON() json.RawMessage
	GameEventType() GameEventType
}

func ToJson(event SerializableEvent) []byte {
	rawEventData := event.RawJSON()
	eventType := event.GameEventType()
	if rawEventData == nil {
		return nil
	}

	gameEvent := GameEvent{
		Type: eventType,
		Data: rawEventData,
	}

	bytes, err := json.Marshal(&gameEvent)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal GameEvent<" + eventType.String() + "> into JSON")
		return nil
	}
	return bytes
}
