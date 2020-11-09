package game

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/dchest/uniuri"
	"github.com/go-chi/chi"
	"github.com/kvnxiao/pictorio/ctxs"
	"github.com/kvnxiao/pictorio/httpw"
	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

const (
	idLength = 9
)

var (
	roomIDRegex = regexp.MustCompile("^[a-zA-Z0-9]{9}$")
)

// Middleware is an http middleware that validates whether or not the specified room ID in the url pattern is valid
// according to the roomIDRegex.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		if err := validateID(roomID); err != nil {
			if err := httpw.Json(w, model.RoomResponse{Exists: false}); err != nil {
				log.Err(err).Str("roomID", roomID).Msg("Invalid room ID")
			}
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), ctxs.KeyRoomID, roomID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateID validates the room id according to the roomIDRegex.
func validateID(roomID string) error {
	if !roomIDRegex.MatchString(roomID) {
		return errors.New("invalid room ID")
	}
	return nil
}

// GenerateRoomID wraps the uniuri package to return a string with a constant length of 9 characters, using alphanumeric
// characters including capitalization [a-zA-Z0-9], representing a room ID.
func GenerateRoomID() string {
	return uniuri.NewLen(idLength)
}
