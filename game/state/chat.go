package state

import (
	"github.com/kvnxiao/pictorio/events"
)

func (g *GameStateProcessor) sendChatAll(chatEvent events.ChatEvent) {
	g.chatHistory.Append(chatEvent)
	g.players.BroadcastEvent(events.Chat(chatEvent))
}
