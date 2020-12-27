package service

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/kvnxiao/pictorio/api"
	"github.com/kvnxiao/pictorio/ctxs"
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

func (s *Service) RegisterRoutes() *Service {
	s.router.Post(api.RoomCreate, func(w http.ResponseWriter, r *http.Request) {
		ro := s.hub.NewRoom()
		if err := response.Json(w, model.RoomResponse{RoomID: ro.ID(), Exists: true}, http.StatusOK); err != nil {
			log.Error().Err(err).Msg("Unable to encode JSON response")
		}
	})

	s.router.Post(api.RoomExists, func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var roomReq model.RoomRequest
		err := decoder.Decode(&roomReq)
		if err != nil {
			respErr := response.Json(w, model.RoomResponse{Exists: false}, http.StatusBadRequest)
			if respErr != nil {
				log.Error().Err(respErr).Msg("Unable to encode JSON response")
			}
		}

		// Check if room id exists
		_, ok := s.hub.Room(roomReq.RoomID)
		respErr := response.Json(w, model.RoomResponse{RoomID: roomReq.RoomID, Exists: ok}, http.StatusOK)
		if respErr != nil {
			log.Error().Err(respErr).Msg("Unable to encode JSON response")
		}
	})

	s.router.Route(api.Room, func(r chi.Router) {
		r.Route("/{roomID}", func(r chi.Router) {
			r.Use(s.roomIDMiddleware)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				roomID, ok := ctxs.RoomID(ctx)
				if !ok {
					if err := response.Json(w, model.RoomResponse{Exists: false}, http.StatusBadRequest); err != nil {
						log.Error().Err(err).Str("route", "/room/"+roomID).Msg("Unable to encode JSON response")
					}
					return
				}
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
				ro.ConnectionHandler(w, r)
			})
		})
	})
	return s
}

func (s *Service) Serve(addr string) {
	log.Info().Msg("Starting server on " + addr)
	err := http.ListenAndServe(addr, s.router)
	if err != nil {
		log.Fatal().Err(err)
	}
}
