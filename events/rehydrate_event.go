package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type RehydrateEvent struct {
	SelfUser        model.User       `json:"selfUser"`
	GameStatus      model.GameStatus `json:"gameStatus"`
	CurrentUserTurn *model.User      `json:"currentUserTurn,omitempty"`
	Lines           []model.Line     `json:"lines"`
}

func RehydrateUser(user model.User) []byte {
	event := RehydrateEvent{
		SelfUser: user,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal RehydrateEvent into JSON.")
		return nil
	}

	return gameEvent(EventTypeRehydrate, eventBytes)
}
