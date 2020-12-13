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

func ReadyUser(user model.User, ready bool) []byte {
	event := ReadyEvent{
		User:  user,
		Ready: ready,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal " + EventTypeReady.String() + " into JSON.")
		return nil
	}

	return gameEvent(EventTypeReady, eventBytes)
}
