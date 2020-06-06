package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kvnxiao/pictorio/ctxs"
	"github.com/kvnxiao/pictorio/room"
	"github.com/rs/zerolog/log"
)

func main() {
	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
	)

	r.Get("/create", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(room.GenerateID()))
	})

	r.Route("/room", func(r chi.Router) {
		r.Route("/{roomID}", func(r chi.Router) {
			r.Use(room.Middleware)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				roomID, ok := ctxs.RoomID(ctx)
				if !ok {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				_, _ = w.Write([]byte(roomID))
			})
		})
	})

	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal().Err(err)
	}
}
