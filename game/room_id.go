package game

import (
	"github.com/dchest/uniuri"
)

// GenerateRoomID wraps the uniuri package to return a string with a constant length of 9 characters, using alphanumeric
// characters including capitalization [a-zA-Z0-9], representing a room ID.
func GenerateRoomID() string {
	return uniuri.NewLen(roomIDLength)
}
