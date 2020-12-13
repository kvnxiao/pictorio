package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type StartGameEvent struct {
	GameStarted bool `json:"gameStarted"`
}

type StartGameIssuedEvent struct {
	Issuer model.User `json:"issuer"`
}

func StartGame() []byte {
	event := StartGameEvent{
		GameStarted: true,
	}
	eventBytes, err := json.Marshal(&event)
	if err != nil {
		log.Err(err).Msg("Could not marshal " + EventTypeStartGame.String() + " into JSON.")
		return nil
	}

	return gameEvent(EventTypeStartGame, eventBytes)
}
