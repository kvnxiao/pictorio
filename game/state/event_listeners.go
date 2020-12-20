package state

import (
	"time"

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

func (g *GameStateProcessor) onChatEvent(event events.ChatEvent) {
	// Do not process chat events from client trying to impersonate the server
	if event.User.ID == model.SystemUserID || event.IsSystem {
		log.Error().
			Msg("Received a " + event.GameEventType().String() + " from client with system user ID / name!")
		return
	}

	// Check if game is in progress and send to guess if so
	if g.status.Status() == model.GameStarted && g.status.TurnStatus() == model.TurnDrawing {
		g.wordGuess <- Guess{
			User:      event.User,
			Timestamp: time.Now().UnixNano(),
			Value:     event.Message,
		}
	} else {
		// Save event to chat history and broadcast
		g.broadcastChat(events.ChatUserEvent(event.User, event.Message))
	}
}

func (g *GameStateProcessor) onDrawEvent(event events.DrawEvent) {
	log.Info().Msg("received draw event!")
	// Validate drawing is from current turn's user
	if g.status.CurrentTurnID() != event.User.ID {
		log.Error().Msg("Received a " + event.GameEventType().String() +
			" from a client whose user ID does not match the current turn's ID")
	}

	// Save event to drawing history
	handled := false
	switch event.Type {
	case events.Line:
		if event.Line == nil {
			log.Error().Msg("Received a " + event.GameEventType().String() + "[Line] event but the line was nil")
		}
		handled = g.drawing.Append(*event.Line)
	case events.Clear:
		handled = g.drawing.Clear()
	case events.Undo:
		handled = g.drawing.Undo()
	case events.Redo:
		handled = g.drawing.Redo()
	default:
		log.Error().Msg("Unknown " + event.GameEventType().String() + " event type")
	}

	// Broadcast event to users
	if handled {
		g.broadcastExcluding(event, event.User.ID)
	}
}

func (g *GameStateProcessor) onReadyEvent(event events.ReadyEvent) {
	ready := g.players.ReadyPlayer(event.User.ID, event.Ready)
	g.broadcast(events.ReadyEvent{User: event.User, Ready: ready})
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
	if g.status.Status() == model.GameStarted {
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
