package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type ChatEvent struct {
	User     model.User `json:"user"`
	Message  string     `json:"message"`
	IsSystem bool       `json:"isSystem,omitempty"`
}

func ChatSystem(message string) []byte {
	event := ChatEvent{
		User:     model.SystemUser(),
		Message:  message,
		IsSystem: true,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal RehydrateEvent into JSON.")
		return nil
	}

	return gameEvent(EventTypeChat, eventBytes)
}

func ChatUser(user model.User, message string) []byte {
	event := ChatEvent{
		User:    user,
		Message: message,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal RehydrateEvent into JSON.")
		return nil
	}

	return gameEvent(EventTypeChat, eventBytes)
}
