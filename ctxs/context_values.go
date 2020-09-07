package ctxs

import (
	"context"
)

type key int

const (
	KeyRoomID key = iota
)

// RoomID gets the roomID string value from the context.
func RoomID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeyRoomID).(string)
	return v, ok
}
