package state

import (
	"context"

	"github.com/kvnxiao/pictorio/events"
	"github.com/kvnxiao/pictorio/game/user"
	"github.com/rs/zerolog/log"
)

func (g *GameStateProcessor) broadcastEvent(eventBytes []byte) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, playerState := range g.playerStates {
		playerState.SendMessage(eventBytes)
	}
}

func (g *GameStateProcessor) sendEvent(userID string, eventBytes []byte) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	playerState, ok := g.playerStates[userID]
	if !ok {
		log.Error().Msg("Attempted to send an event to an invalid player ID")
	}
	playerState.SendMessage(eventBytes)
}

func (g *GameStateProcessor) saveUserConnection(user *user.User) PlayerState {
	g.mu.Lock()
	defer g.mu.Unlock()

	playerState, ok := g.playerStates[user.ID]
	if ok {
		// Existing user re-joined the room
		playerState.SetNewConnection(user)
	} else {
		// No player state exists for this user, create a new player state for the current game
		playerState = NewPlayer(user, g.IsFull())
		g.playerStates[user.ID] = playerState
	}
	playerState.SetConnected(true)
	return playerState
}

func (g *GameStateProcessor) removeUserConnection(userID string) PlayerState {
	playerState, ok := g.playerStates[userID]
	if !ok {
		return nil
	}
	playerState.SetConnected(false)
	return playerState
}

func (g *GameStateProcessor) UserJoined(ctx context.Context, user *user.User, connErrChan chan error) {
	// save user connection
	playerState := g.saveUserConnection(user)

	// concurrently handle the user's WebSocket connection
	go user.ReaderLoop(ctx, g.messageQueue, connErrChan)
	go user.WriterLoop(ctx, connErrChan)

	userModel := playerState.UserModel()

	// send rehydration event to user who just joined
	playerState.SendMessage(events.RehydrateUser(userModel))

	// broadcast that a user has joined
	g.broadcastEvent(events.UserJoin(userModel))
}

func (g *GameStateProcessor) UserLeft(userID string) {
	playerState := g.removeUserConnection(userID)

	// broadcast that a user has left
	if playerState != nil {
		g.broadcastEvent(events.UserLeave(playerState.UserModel()))
	}
}
