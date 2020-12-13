package state

import (
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

func (g *GameStateProcessor) onChatEvent(chatEvent events.ChatEvent) {
	// Do not process chat events from client trying to impersonate the server
	if chatEvent.User.ID == model.SystemUserID || chatEvent.IsSystem {
		log.Error().
			Msg("Received a ChatEvent from client with system user ID / name!")
		return
	}

	// save event to chat history
	g.chatHistory.Append(chatEvent)
	g.broadcastEvent(events.Chat(chatEvent))
}

func (g *GameStateProcessor) onDrawEvent(event events.GameEvent) {
	log.Info().Msg("onDrawEvent")
}
