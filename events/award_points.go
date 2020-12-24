package events

import (
	"encoding/json"

	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

type AwardPointsEvent struct {
	Guesser       model.User `json:"guesser"`
	GuesserPoints int        `json:"guesserPoints"`
	Drawer        model.User `json:"drawer"`
	DrawerPoints  int        `json:"drawerPoints"`
}

func (e AwardPointsEvent) RawJSON() json.RawMessage {
	eventBytes, err := json.Marshal(e)
	if err != nil {
		log.Error().Err(err).Msg("Could not marshal " + e.GameEventType().String() + " into JSON.")
		return nil
	}
	return eventBytes
}

func (e AwardPointsEvent) GameEventType() GameEventType {
	return EventTypeAwardPoints
}
