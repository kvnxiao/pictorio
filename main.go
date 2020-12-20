package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/kvnxiao/pictorio/service"
	"github.com/rs/zerolog"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var hostFlag = flag.String("host", ":3000", "The hostname to start the server on")
	var debugFlag = flag.Bool("debug", false, "Enables debug mode for logging")

	flag.Parse()

	if *debugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	server := service.NewService()
	server.
		SetupMiddleware().
		RegisterRoutes().
		Serve(*hostFlag)
}
