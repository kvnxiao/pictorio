package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type GameEventType int

const (
	EventTypePlayerAction GameEventType = iota
	EventTypeSelfJoin
	EventTypeRehydrate
	EventTypeChat
	EventTypeDraw
)

type PlayerAction int

const (
	PlayerActionLeave PlayerAction = iota
	PlayerActionJoin
)

type SelfJoinEvent struct {
	Player model.Player `json:"player"`
}

type PlayerActionEvent struct {
	Player model.Player `json:"player"`
	Action PlayerAction `json:"action"`
}

type GameEvent struct {
	Type GameEventType `json:"type"`
	Data interface{}   `json:"data"`
}

func SelfJoinEventMessage(player model.Player) []byte {
	event := SelfJoinEvent{
		Player: player,
	}
	gameEvent := GameEvent{
		Type: EventTypeSelfJoin,
		Data: event,
	}
	bytes, err := json.Marshal(&gameEvent)
	if err != nil {
		log.Err(err).Msg("Could not marshal SelfJoinEvent into JSON.")
		return nil
	}
	return bytes
}
