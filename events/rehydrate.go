package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type RehydrateEvent struct {
	SelfUser        model.User             `json:"selfUser"`
	CurrentTurnUser *model.User            `json:"currentTurnUser"`
	ChatMessages    []ChatEvent            `json:"chatMessages"`
	Players         model.PlayersSummary   `json:"players"`
	Game            model.GameStateSummary `json:"game"`
	Lines           []model.Line           `json:"lines"`
}

func (e RehydrateEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e RehydrateEvent) GameEventType() GameEventType {
	return EventTypeRehydrate
}

func RehydrateForUser(
	selfUser model.User,
	currentTurnUser *model.User,
	chatHistory []ChatEvent,
	playersSummary model.PlayersSummary,
	gameSummary model.GameStateSummary,
	lines []model.Line,
) RehydrateEvent {
	return RehydrateEvent{
		SelfUser:        selfUser,
		CurrentTurnUser: currentTurnUser,
		ChatMessages:    chatHistory,
		Players:         playersSummary,
		Game:            gameSummary,
		Lines:           lines,
	}
}
