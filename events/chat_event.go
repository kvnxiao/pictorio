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

func ChatSystemEvent(message string) ChatEvent {
	return ChatEvent{
		User:     model.SystemUser(),
		Message:  message,
		IsSystem: true,
	}
}

func ChatUserEvent(user model.User, message string) ChatEvent {
	return ChatEvent{
		User:    user,
		Message: message,
	}
}

func Chat(event ChatEvent) []byte {
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal ChatEvent into JSON.")
		return nil
	}

	return gameEvent(EventTypeChat, eventBytes)
}
