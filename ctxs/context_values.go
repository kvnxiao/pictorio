package ctxs

import (
	"context"

	"github.com/segmentio/ksuid"
)

type key int

const (
	KeyRoomID key = iota
	KeyUserID
	KeyUserName
)

// RoomID gets the roomID string value from the context.
func RoomID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeyRoomID).(string)
	return v, ok
}

// UserName gets the name of the user from the context.
func UserName(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(KeyUserName).(string)
	return v, ok
}

// UserID gets the ID of the user from the context.
func UserID(ctx context.Context) (ksuid.KSUID, bool) {
	v, ok := ctx.Value(KeyUserID).(ksuid.KSUID)
	return v, ok
}
