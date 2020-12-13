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
	PlayerState model.PlayerState   `json:"playerState"`
	Action      UserJoinLeaveAction `json:"action"`
}

func joinLeaveEvent(playerState model.PlayerState, action UserJoinLeaveAction) []byte {
	event := UserJoinLeaveEvent{
		PlayerState: playerState,
		Action:      action,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal " + EventTypeUserJoinLeave.String() + " into JSON")
		return nil
	}

	return gameEvent(EventTypeUserJoinLeave, eventBytes)
}

func UserJoin(playerState model.PlayerState) []byte {
	return joinLeaveEvent(playerState, UserActionJoin)
}

func UserLeave(playerState model.PlayerState) []byte {
	return joinLeaveEvent(playerState, UserActionLeave)
}
