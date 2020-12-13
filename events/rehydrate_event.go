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
	CurrentUserTurn *model.User      `json:"currentUserTurn,omitempty"`
	Lines           []model.Line     `json:"lines"`
}

func RehydrateForUser(
	user model.User,
	playerStates []model.PlayerState,
	chatHistory []ChatEvent,
	status model.GameStatus,
	maxPlayerCount int,
) []byte {
	event := RehydrateEvent{
		SelfUser:     user,
		PlayerStates: playerStates,
		ChatMessages: chatHistory,
		GameStatus:   status,
		MaxPlayers:   maxPlayerCount,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal " + EventTypeRehydrate.String() + " into JSON.")
		return nil
	}

	return gameEvent(EventTypeRehydrate, eventBytes)
}
