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

func (e ChatEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e ChatEvent) GameEventType() GameEventType {
	return EventTypeChat
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
