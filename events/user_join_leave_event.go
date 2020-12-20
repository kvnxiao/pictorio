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

func (e UserJoinLeaveEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e UserJoinLeaveEvent) GameEventType() GameEventType {
	return EventTypeUserJoinLeave
}

func joinLeaveEvent(playerState model.PlayerState, action UserJoinLeaveAction) UserJoinLeaveEvent {
	return UserJoinLeaveEvent{
		PlayerState: playerState,
		Action:      action,
	}
}

func UserJoin(playerState model.PlayerState) UserJoinLeaveEvent {
	return joinLeaveEvent(playerState, UserActionJoin)
}

func UserLeave(playerState model.PlayerState) UserJoinLeaveEvent {
	return joinLeaveEvent(playerState, UserActionLeave)
}
