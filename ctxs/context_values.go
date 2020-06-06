package ctxs

import (
	"context"
)

// RoomID gets the roomID string value from the context.
func RoomID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value("roomID").(string)
	return v, ok
}
