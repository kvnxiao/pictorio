package events

type GameEventType int

const (
	EventTypePlayerJoinLeaveEvent GameEventType = iota
	EventTypeSelfJoin
	EventTypeRehydrate
	EventTypeChat
	EventTypeDraw
)

type GameEvent struct {
	Type GameEventType `json:"type"`
	Data interface{}   `json:"data"`
}
