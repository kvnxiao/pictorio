package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type TurnNextPlayerNonce struct {
	NextTurnUser model.User `json:"nextTurnUser"`
}

type TurnNextPlayerEvent struct {
	Nonce    *TurnNextPlayerNonce `json:"nonce,omitempty"`
	MaxTime  int                  `json:"maxTime"`
	TimeLeft int                  `json:"timeLeft"`
	Status   model.TurnStatus     `json:"status"`
}

func (e TurnNextPlayerEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e TurnNextPlayerEvent) GameEventType() GameEventType {
	return EventTypeTurnNextPlayer
}

func TurnBeginNextPlayer(nextPlayer model.User, maxTimeSeconds int) TurnNextPlayerEvent {
	return TurnNextPlayerEvent{
		Nonce: &TurnNextPlayerNonce{
			NextTurnUser: nextPlayer,
		},
		MaxTime:  maxTimeSeconds,
		TimeLeft: maxTimeSeconds,
		Status:   model.TurnNextPlayer,
	}
}

func TurnNextPlayerCountdown(maxTimeSeconds int, timeLeftSeconds int) TurnNextPlayerEvent {
	return TurnNextPlayerEvent{
		Nonce:    nil,
		MaxTime:  maxTimeSeconds,
		TimeLeft: timeLeftSeconds,
		Status:   model.TurnNextPlayer,
	}
}
