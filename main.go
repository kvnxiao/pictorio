package main

import (
	"flag"

	"github.com/kvnxiao/pictorio/service"
)

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
