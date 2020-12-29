package state

import (
	"time"

	"github.com/kvnxiao/pictorio/events"
	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

func (g *GameStateProcessor) warnServerSourcedEvent(eventType events.GameEventType) {
	log.Warn().
		Str("event", eventType.String()).
		Msg("Received a server-sourced event from a client!")
}

func (g *GameStateProcessor) onChatEvent(event events.ChatEvent) {
	// Do not process chat events from client trying to impersonate the server
	if event.User.ID == model.SystemUserID || event.Type != events.ChatEventUser {
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
		g.broadcastChat(events.ChatUserMessage(event.User, event.Message))
	}
}

func (g *GameStateProcessor) onDrawEvent(event events.DrawEvent) {
	// Validate drawing is from current turn's user
	if g.status.CurrentTurnID() != event.User.ID {
		log.Error().Msg("Received a " + event.GameEventType().String() +
			" from a client whose user ID does not match the current turn's ID")
	}

	// Save event to drawing history
	handled := false
	switch event.Type {
	case events.Line:
		handled = g.drawingHistory.PromoteLine()
	case events.Clear:
		handled = g.drawingHistory.Clear()
	case events.Undo:
		handled = g.drawingHistory.Undo()
	case events.Redo:
		handled = g.drawingHistory.Redo()
	default:
		log.Error().Msg("Unknown " + event.GameEventType().String() + " event type")
	}

	// Broadcast event to users
	if handled {
		g.broadcastExcluding(event, event.User.ID)
	}
}

func (g *GameStateProcessor) onDrawTempEvent(event events.DrawTempEvent) {
	// Validate drawing is from current turn's user
	if g.status.CurrentTurnID() != event.User.ID {
		log.Error().Msg("Received a " + event.GameEventType().String() +
			" from a client whose user ID does not match the current turn's ID")
	}

	g.drawingHistory.AppendFromTempLine(event.Line)
	g.broadcastExcluding(event, event.User.ID)
}

func (g *GameStateProcessor) onDrawSelectColour(event events.DrawSelectColourEvent) {
	// Validate drawing is from current turn's user
	if g.status.CurrentTurnID() != event.User.ID {
		log.Error().Msg("Received a " + event.GameEventType().String() +
			" from a client whose user ID does not match the current turn's ID")
	}

	g.drawingHistory.SetTempColour(event.ColourIndex)
	g.broadcastExcluding(event, event.User.ID)
}

func (g *GameStateProcessor) onDrawSelectThickness(event events.DrawSelectThicknessEvent) {
	// Validate drawing is from current turn's user
	if g.status.CurrentTurnID() != event.User.ID {
		log.Error().Msg("Received a " + event.GameEventType().String() +
			" from a client whose user ID does not match the current turn's ID")
	}

	g.drawingHistory.SetTempThickness(event.ThicknessIndex)
	g.broadcastExcluding(event, event.User.ID)
}

func (g *GameStateProcessor) onReadyEvent(event events.ReadyEvent) {
	ready := g.players.ReadyPlayer(event.User.ID, event.Ready)
	g.broadcast(events.ReadyEvent{User: event.User, Ready: ready})
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
	log.Debug().Msg("Game started!")
}

func (g *GameStateProcessor) onTurnWordSelectedEvent(event events.TurnWordSelectedEvent) {
	// Validate the user who sent this event is the current turn's user
	if event.User.ID != g.status.CurrentTurnID() {
		log.Error().Msg("Received a " + events.EventTypeTurnWordSelected.String() +
			" event from a player who's turn is not the current turn.")
		return
	}

	g.wordSelectionIndex <- SelectionIndex{
		User:      event.User,
		Timestamp: time.Now().UnixNano(),
		Value:     event.Index,
	}
}

func (g *GameStateProcessor) onNewGameIssued(event events.NewGameIssuedEvent) {
	// Validate the issuer is the room leader
	if event.Issuer.ID != g.players.RoomLeaderID() {
		log.Error().
			Msg("Received a " + events.EventTypeNewGameIssued.String() +
				" from client who was not the room leader!")
		return
	}

	g.status.Reset()
	g.players.Reset()
	g.status.SetStatus(model.GameWaitingReadyUp)
	g.broadcast(events.NewGameResetEvent{
		PlayerStates: g.players.Summary().PlayerStates,
	})
}
