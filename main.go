package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kvnxiao/pictorio/ctxs"
	"github.com/kvnxiao/pictorio/gameserver"
	"github.com/kvnxiao/pictorio/room"
	"github.com/rs/zerolog/log"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
	)

	gs := gameserver.New()

	r.HandleFunc("/", indexHandler)

	r.Get("/create", func(w http.ResponseWriter, r *http.Request) {
		ro := gs.NewRoom()
		http.Redirect(w, r, "/room/"+ro.ID(), http.StatusSeeOther)
		log.Info().Str("roomID", ro.ID()).Send()
	})

	r.Route("/room", func(r chi.Router) {
		r.Route("/{roomID}", func(r chi.Router) {
			r.Use(room.Middleware)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				_, ok := ctxs.RoomID(ctx)
				if !ok {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				indexHandler(w, r)
			})
			r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				roomID, ok := ctxs.RoomID(ctx)
				if !ok {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				ro, ok := gs.Room(roomID)
				if !ok {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				log.Info().Str("roomID", roomID).Send()
				ro.ConnectionHandler(w, r)
			})
		})
	})

	log.Info().Msg("Serving on :3000")
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal().Err(err)
	}
}
