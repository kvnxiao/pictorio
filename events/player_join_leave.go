package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type PlayerJoinLeaveAction int

const (
	PlayerActionJoin PlayerJoinLeaveAction = iota
	PlayerActionLeave
)

type PlayerJoinLeaveEvent struct {
	Player model.User            `json:"player"`
	Action PlayerJoinLeaveAction `json:"action"`
}

func joinLeaveEvent(player model.User, action PlayerJoinLeaveAction) []byte {
	event := PlayerJoinLeaveEvent{
		Player: player,
		Action: action,
	}
	gameEvent := GameEvent{
		Type: EventTypePlayerJoinLeaveEvent,
		Data: event,
	}
	bytes, err := json.Marshal(&gameEvent)
	if err != nil {
		log.Err(err).Msg("Could not marshal PlayerJoinLeaveEvent into JSON.")
		return nil
	}
	return bytes
}

func PlayerJoin(player model.User) []byte {
	return joinLeaveEvent(player, PlayerActionJoin)
}

func PlayerLeave(player model.User) []byte {
	return joinLeaveEvent(player, PlayerActionLeave)
}
