package events

type GameEventType int

const (
	EventTypeUserJoinLeaveEvent GameEventType = iota
	EventTypeRehydrate
	EventTypeChat
	EventTypeDraw
)

type GameEvent struct {
	Type GameEventType `json:"type"`
	Data interface{}   `json:"data"`
}
