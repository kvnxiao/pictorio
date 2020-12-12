package events

import (
	"encoding/json"
)

type GameEventType int

const (
	EventTypeUserJoinLeaveEvent GameEventType = iota
	EventTypeRehydrate
	EventTypeChat
	EventTypeDraw
)

type GameEvent struct {
	Type GameEventType   `json:"type"`
	Data json.RawMessage `json:"data"`
}
