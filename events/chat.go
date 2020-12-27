package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type ChatEventType int

const (
	ChatEventSystem ChatEventType = iota
	ChatEventUser
	ChatEventJoin
	ChatEventLeave
	ChatEventGuessed
)

const (
	userJoinedMsg = "has joined the room."
	userLeftMsg   = "has left the room."
	userGuessed   = "has guessed the word."

	formatSystem     = "%m"
	formatUser       = "%u: %m"
	formatUserAction = "%u %m"
)

type ChatEvent struct {
	User    model.User    `json:"user"`
	Message string        `json:"message"`
	Format  string        `json:"format"`
	Type    ChatEventType `json:"type"`
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

func ChatUserJoined(user model.User) ChatEvent {
	return ChatEvent{
		User:    user,
		Message: userJoinedMsg,
		Format:  formatUserAction,
		Type:    ChatEventJoin,
	}
}

func ChatUserLeft(user model.User) ChatEvent {
	return ChatEvent{
		User:    user,
		Message: userLeftMsg,
		Format:  formatUserAction,
		Type:    ChatEventLeave,
	}
}

func ChatUserGuessed(user model.User) ChatEvent {
	return ChatEvent{
		User:    user,
		Message: userGuessed,
		Format:  formatUserAction,
		Type:    ChatEventGuessed,
	}
}

func ChatSystemEvent(message string) ChatEvent {
	return ChatEvent{
		User:    model.SystemUser(),
		Message: message,
		Format:  formatSystem,
		Type:    ChatEventSystem,
	}
}

func ChatUserMessage(user model.User, message string) ChatEvent {
	return ChatEvent{
		User:    user,
		Message: message,
		Format:  formatUser,
		Type:    ChatEventUser,
	}
}
