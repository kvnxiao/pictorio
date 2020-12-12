package state

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/events"
	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

func (g *GameStateProcessor) onUserJoinLeaveEvent() {
	log.Warn().Msg("Received a UserJoinLeaveEvent when this is supposed to be a client only event!")
}

func (g *GameStateProcessor) onRehydrateEvent() {
	log.Warn().Msg("Received a RehydrateEvent when this is supposed to be a client only event!")
}

func (g *GameStateProcessor) onChatEvent(event events.GameEvent) {
	var chatEvent model.ChatEvent
	err := json.Unmarshal(event.Data, &chatEvent)
	if err != nil {
		log.Error().
			Err(err).Msg("Could not unmarshal ChatEvent from user")
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Error().
			Err(err).Msg("Could not send ChatEvent to user")
	}

	g.broadcastEvent(eventBytes)
}

func (g *GameStateProcessor) onDrawEvent(event events.GameEvent) {
	log.Info().Msg("onDrawEvent")
}
