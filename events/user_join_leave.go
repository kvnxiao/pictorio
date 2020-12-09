package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type UserJoinLeaveAction int

const (
	UserActionJoin UserJoinLeaveAction = iota
	UserActionLeave
)

type UserJoinLeaveEvent struct {
	User   model.User          `json:"user"`
	Action UserJoinLeaveAction `json:"action"`
}

func joinLeaveEvent(user model.User, action UserJoinLeaveAction) []byte {
	event := UserJoinLeaveEvent{
		User:   user,
		Action: action,
	}
	gameEvent := GameEvent{
		Type: EventTypeUserJoinLeaveEvent,
		Data: event,
	}
	bytes, err := json.Marshal(&gameEvent)
	if err != nil {
		log.Err(err).Msg("Could not marshal UserJoinLeaveEvent into JSON.")
		return nil
	}
	return bytes
}

func UserJoin(user model.User) []byte {
	return joinLeaveEvent(user, UserActionJoin)
}

func UserLeave(user model.User) []byte {
	return joinLeaveEvent(user, UserActionLeave)
}
