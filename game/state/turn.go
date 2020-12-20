package state

import (
	"errors"

	"github.com/kvnxiao/pictorio/model"
)

type SelectionIndex struct {
	User      model.User
	Timestamp int64
	Value     int
}

type Guess struct {
	User      model.User
	Timestamp int64
	Value     string
}

func (g *GameStateProcessor) getCurrentTurnUser() (model.User, error) {
	currentTurnID := g.status.CurrentTurnID()
	player, ok := g.players.GetPlayer(currentTurnID)
	if !ok {
		return model.User{}, errors.New("could not get player state with invalid player id")
	}

	userModel := player.ToUserModel()
	return userModel, nil
}
