package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type RehydrateEvent struct {
	// UserRehydrateEvent
	SelfUser     model.User          `json:"selfUser"`
	PlayerStates []model.PlayerState `json:"playerStates"`

	// ChatRehydrateEvent
	ChatMessages []ChatEvent `json:"chatMessages"`

	// GameRehydrateEvent
	MaxPlayers      int              `json:"maxPlayers"`
	GameStatus      model.GameStatus `json:"gameStatus"`
	PlayerOrderIDs  []string         `json:"playerOrderIds"`
	CurrentUserTurn *model.User      `json:"currentUserTurn"`
	Lines           []model.Line     `json:"lines,omitempty"`
}

func RehydrateForUser(
	user model.User,
	playerStates []model.PlayerState,
	chatHistory []ChatEvent,
	status model.GameStatus,
	maxPlayerCount int,
	playerOrderIDs []string,
	currentTurnUser *model.User,
	lines []model.Line,
) []byte {
	event := RehydrateEvent{
		SelfUser:        user,
		PlayerStates:    playerStates,
		ChatMessages:    chatHistory,
		GameStatus:      status,
		MaxPlayers:      maxPlayerCount,
		PlayerOrderIDs:  playerOrderIDs,
		CurrentUserTurn: currentTurnUser,
		Lines:           lines,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal " + EventTypeRehydrate.String() + " into JSON.")
		return nil
	}

	return gameEvent(EventTypeRehydrate, eventBytes)
}
