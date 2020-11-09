package service

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kvnxiao/pictorio/ctxs"
	"github.com/kvnxiao/pictorio/game"
	"github.com/kvnxiao/pictorio/hub"
	"github.com/kvnxiao/pictorio/model"
	"github.com/kvnxiao/pictorio/response"
	"github.com/rs/zerolog/log"
)

type Service struct {
	hub    *hub.Hub
	router chi.Router
}

func NewService() *Service {
	return &Service{
		hub:    hub.New(),
		router: chi.NewRouter(),
	}
}

func (s *Service) SetupMiddleware() *Service {
	s.router.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
	)
	return s
}

func (s *Service) FileServer() *Service {
	for _, file := range files {
		handleFolder(s.router, file)
	}
	return s
}

func (s *Service) RegisterRoutes() *Service {
	s.router.NotFound(indexHandler)

	s.router.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		ro := s.hub.NewRoom()
		if err := response.Json(w, model.RoomResponse{RoomID: ro.ID(), Exists: true}); err != nil {
			log.Err(err).Msg("Unable to encode JSON response")
		}
	})

	s.router.Route("/room", func(r chi.Router) {
		r.Route("/{roomID}", func(r chi.Router) {
			r.Use(game.Middleware)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				roomID, ok := ctxs.RoomID(ctx)
				if !ok {
					if err := response.Json(w, model.RoomResponse{Exists: false}); err != nil {
						log.Err(err).Str("route", "/room/"+roomID).Msg("Unable to encode JSON response")
					}
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
				ro, ok := s.hub.Room(roomID)
				if !ok {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				log.Info().Str("roomID", roomID).Send()
				ro.ConnectionHandler(w, r)
			})
		})
	})
	return s
}

func (s *Service) Serve() {
	port := ":3000"
	log.Info().Msg("Starting server on port " + port)
	err := http.ListenAndServe(port, s.router)
	if err != nil {
		log.Fatal().Err(err)
	}
}
