package events

import (
	"github.com/kvnxiao/pictorio/model"
)

// TurnWordSelectedEvent is the client-sourced event in which the current turn player selects a word from the generated
// list of words that the server produces (this list is sent to the client from beforehand in a TurnEvent)
type TurnWordSelectedEvent struct {
	User  model.User `json:"user"`
	Index int        `json:"index"`
}
