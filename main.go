package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/kvnxiao/pictorio/ctxs"
	"github.com/kvnxiao/pictorio/fs"
	"github.com/kvnxiao/pictorio/game"
	"github.com/kvnxiao/pictorio/httpw"
	"github.com/kvnxiao/pictorio/hub"
	"github.com/kvnxiao/pictorio/model"
	"github.com/rs/zerolog/log"
)

func main() {
	http.FileServer(http.Dir("./dist"))
	r := chi.NewRouter()
	r.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"http://localhost:8080"},
		}),
	)

	fs.FileServer(r)
	r.NotFound(fs.IndexHandler)

	gs := hub.New()

	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		ro := gs.NewRoom()
		if err := httpw.Json(w, model.RoomResponse{RoomID: ro.ID(), Exists: true}); err != nil {
			log.Err(err).Msg("Unable to encode JSON response")
		}
	})

	r.Route("/room", func(r chi.Router) {
		r.Route("/{roomID}", func(r chi.Router) {
			r.Use(game.Middleware)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				roomID, ok := ctxs.RoomID(ctx)
				if !ok {
					if err := httpw.Json(w, model.RoomResponse{Exists: false}); err != nil {
						log.Err(err).Str("route", "/room/"+roomID).Msg("Unable to encode JSON response")
					}
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				fs.IndexHandler(w, r)
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
