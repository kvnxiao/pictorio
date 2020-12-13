package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/kvnxiao/pictorio/service"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.Parse()

	addr := flag.Arg(0)
	if addr == "" {
		addr = ":3000"
	}

	server := service.NewService()
	server.
		SetupMiddleware().
		FileServer().
		RegisterRoutes().
		Serve(addr)
}
