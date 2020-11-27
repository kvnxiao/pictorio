package ctxs

import (
	"context"

	"github.com/segmentio/ksuid"
)

type key int

const (
	KeyRoomID key = iota
	KeyPlayerID
	KeyPlayerName
)

// RoomID gets the roomID string value from the context.
func RoomID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeyRoomID).(string)
	return v, ok
}

// PlayerName gets the name of the player from the context.
func PlayerName(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeyPlayerName).(string)
	return v, ok
}

// PlayerID gets the ID of the player from the context.
func PlayerID(ctx context.Context) (ksuid.KSUID, bool) {
	v, ok := ctx.Value(KeyPlayerID).(ksuid.KSUID)
	return v, ok
}
