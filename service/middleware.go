package service

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/go-chi/chi"
	"github.com/kvnxiao/pictorio/ctxs"
)

var (
	roomIDRegex = regexp.MustCompile("^[a-zA-Z0-9]{9}$")
)

// roomIdMiddleware is an http middleware that validates whether or not the specified room ID in the url pattern is
// valid according to the roomIDRegex, and whether the room exists.
func (s *Service) roomIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		if err := validateID(roomID); err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		_, ok := s.hub.Room(roomID)
		if !ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
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
