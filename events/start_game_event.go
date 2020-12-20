package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type StartGameEvent struct {
	PlayerOrderIDs  []string   `json:"playerOrderIds"`
	CurrentUserTurn model.User `json:"currentUserTurn"`
}

type StartGameIssuedEvent struct {
	Issuer model.User `json:"issuer"`
}

func StartGame(playerOrderIDs []string, currentUserTurn model.User) []byte {
	event := StartGameEvent{
		PlayerOrderIDs:  playerOrderIDs,
		CurrentUserTurn: currentUserTurn,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal " + EventTypeStartGame.String() + " into JSON.")
		return nil
	}

	return gameEvent(EventTypeStartGame, eventBytes)
}
