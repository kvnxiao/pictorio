package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/kvnxiao/pictorio/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var hostFlag = flag.String("host", ":3000", "The hostname to start the server on")
	var debugFlag = flag.Bool("debug", false, "Enables debug mode for logging")

	flag.Parse()

	if *debugFlag {
		log.Level(zerolog.DebugLevel)
	}

	server := service.NewService()
	server.
		SetupMiddleware().
		FileServer().
		RegisterRoutes().
		Serve(*hostFlag)
}
