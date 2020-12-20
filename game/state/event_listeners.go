package state

import (
	"github.com/kvnxiao/pictorio/events"
	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

func (g *GameStateProcessor) onUserJoinLeaveEvent() {
	log.Warn().
		Str("event", events.EventTypeUserJoinLeave.String()).
		Msg("Received a server-sourced event from a client!")
}

func (g *GameStateProcessor) onRehydrateEvent() {
	log.Warn().
		Str("event", events.EventTypeRehydrate.String()).
		Msg("Received a server-sourced event from a client!")
}

func (g *GameStateProcessor) onChatEvent(chatEvent events.ChatEvent) {
	// Do not process chat events from client trying to impersonate the server
	if chatEvent.User.ID == model.SystemUserID || chatEvent.IsSystem {
		log.Error().
			Msg("Received a " + events.EventTypeChat.String() + " from client with system user ID / name!")
		return
	}

	// save event to chat history
	eventCopy := events.ChatUserEvent(chatEvent.User, chatEvent.Message)
	g.chatHistory.Append(eventCopy)
	g.players.BroadcastEvent(events.Chat(eventCopy))
}

func (g *GameStateProcessor) onDrawEvent(event events.DrawEvent) {
	log.Info().Msg("received draw event!")
	// Validate drawing is from current turn's user
	if g.status.CurrentTurnID() != event.User.ID {
		log.Error().Msg("Received a " + events.EventTypeDraw.String() +
			" from a client whose user ID does not match the current turn's ID")
	}

	// Save event to drawing history
	switch event.Type {
	case events.Line:
		if event.Line == nil {
			log.Error().Msg("Received a " + events.EventTypeDraw.String() + "[Line] event but the line was nil")
		}
		g.drawing.Append(*event.Line)
	case events.Clear:
		g.drawing.Clear()
	case events.Undo:
		g.drawing.Undo()
	case events.Redo:
		g.drawing.Redo()
	default:
		log.Error().Msg("Unknown " + events.EventTypeDraw.String() + " event type")
	}

	// Broadcast event to users
	g.players.BroadcastEventExclude(events.Draw(event), event.User.ID)
}

func (g *GameStateProcessor) onReadyEvent(event events.ReadyEvent) {
	ready := g.players.ReadyPlayer(event.User.ID, event.Ready)
	log.Info().
		Str("uid", event.User.ID).
		Bool("ready", ready).
		Msg("User switched ready state")
	g.players.BroadcastEvent(events.ReadyUser(event.User, ready))
}

func (g *GameStateProcessor) onStartGameEvent() {
	log.Warn().
		Str("event", events.EventTypeStartGame.String()).
		Msg("Received a server-sourced event from a client!")
}

func (g *GameStateProcessor) onStartGameIssuedEvent(event events.StartGameIssuedEvent) {
	// Validate the issuer is the room leader
	if event.Issuer.ID != g.players.RoomLeaderID() {
		log.Error().
			Msg("Received a " + events.EventTypeStartGameIssued.String() +
				" from client who was not the room leader!")
		return
	}

	// Validate that the game has not already started
	if g.status.Status() == model.StatusStarted {
		log.Error().Msg("Received a " + events.EventTypeStartGameIssued.String() +
			" event from the room leader, but the game has already started!")
		return
	}

	started := g.StartGame()
	if !started {
		log.Warn().Msg("Failed to start game!")
	}
	log.Info().Msg("Game started!")
}
