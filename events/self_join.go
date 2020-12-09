package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type RehydrateEvent struct {
	User model.User `json:"user"`
}

func RehydrateUser(user model.User) []byte {
	event := RehydrateEvent{
		User: user,
	}
	gameEvent := GameEvent{
		Type: EventTypeRehydrate,
		Data: event,
	}
	bytes, err := json.Marshal(&gameEvent)
	if err != nil {
		log.Err(err).Msg("Could not marshal RehydrateEvent into JSON.")
		return nil
	}
	return bytes
}
