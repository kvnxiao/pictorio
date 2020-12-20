package state

import (
	"github.com/kvnxiao/pictorio/events"
)

func (g *GameStateProcessor) broadcast(event events.SerializableEvent) {
	g.players.SendEventToAll(event)
}

func (g *GameStateProcessor) broadcastExcluding(event events.SerializableEvent, userID string) {
	g.players.SendEventToAllExcept(event, userID)
}

func (g *GameStateProcessor) emit(event events.SerializableEvent, userID string) {
	g.players.SendEventToUser(event, userID)
}

func (g *GameStateProcessor) broadcastChat(chatEvent events.ChatEvent) {
	g.chatHistory.Append(chatEvent)
	g.players.SendEventToAll(chatEvent)
}
