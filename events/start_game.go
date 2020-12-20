package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type StartGameEvent struct {
	PlayerOrderIDs  []string   `json:"playerOrderIds"`
	CurrentUserTurn model.User `json:"currentUserTurn"`
}

func (e StartGameEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e StartGameEvent) GameEventType() GameEventType {
	return EventTypeStartGame
}

type StartGameIssuedEvent struct {
	Issuer model.User `json:"issuer"`
}
