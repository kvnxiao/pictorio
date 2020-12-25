package events

import (
	"encoding/json"
)

type GameEventType int

const (
	EventTypeUserJoinLeave     GameEventType = iota // server-sourced
	EventTypeRehydrate                              // server-sourced
	EventTypeChat                                   // bi-directional
	EventTypeDraw                                   // bi-directional
	EventTypeReady                                  // bi-directional
	EventTypeStartGame                              // server-sourced
	EventTypeStartGameIssued                        // client-sourced
	EventTypeTurnNextPlayer                         // server-sourced
	EventTypeTurnWordSelection                      // server-sourced
	EventTypeTurnWordSelected                       // client-sourced
	EventTypeTurnDrawing                            // server-sourced
	EventTypeTurnEnd                                // server-sourced
	EventTypeAwardPoints                            // server-sourced
	EventTypeGameOver                               // server-sourced
)

func (e GameEventType) String() string {
	switch e {
	case EventTypeUserJoinLeave:
		return "UserJoinLeaveEvent"
	case EventTypeRehydrate:
		return "RehydrateEvent"
	case EventTypeChat:
		return "ChatEvent"
	case EventTypeDraw:
		return "DrawEvent"
	case EventTypeReady:
		return "ReadyEvent"
	case EventTypeStartGame:
		return "StartGameEvent"
	case EventTypeStartGameIssued:
		return "StartGameIssuedEvent"
	case EventTypeTurnNextPlayer:
		return "TurnNextPlayerEvent"
	case EventTypeTurnWordSelection:
		return "TurnWordSelectionEvent"
	case EventTypeTurnWordSelected:
		return "TurnWordSelectedEvent"
	case EventTypeTurnDrawing:
		return "TurnDrawingEvent"
	case EventTypeTurnEnd:
		return "TurnEndEvent"
	case EventTypeAwardPoints:
		return "AwardPointsEvent"
	case EventTypeGameOver:
		return "GameOverEvent"
	default:
		return "UNKNOWN_Event"
	}
}

type GameEvent struct {
	Type GameEventType   `json:"type"`
	Data json.RawMessage `json:"data"`
}
