package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type SelfJoinEvent struct {
	Player model.User `json:"player"`
}

func SelfJoinEventMessage(player model.User) []byte {
	event := SelfJoinEvent{
		Player: player,
	}
	gameEvent := GameEvent{
		Type: EventTypeSelfJoin,
		Data: event,
	}
	bytes, err := json.Marshal(&gameEvent)
	if err != nil {
		log.Err(err).Msg("Could not marshal SelfJoinEvent into JSON.")
		return nil
	}
	return bytes
}
