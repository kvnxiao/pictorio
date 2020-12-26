package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type NewGameResetEvent struct {
	PlayerStates []model.PlayerState `json:"playerStates"`
}

func (e NewGameResetEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e NewGameResetEvent) GameEventType() GameEventType {
	return EventTypeNewGameReset
}

type NewGameIssuedEvent struct {
	Issuer model.User `json:"issuer"`
}
