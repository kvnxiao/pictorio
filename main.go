package main

import (
	"github.com/kvnxiao/pictorio/service"
)

func main() {
	server := service.NewService()
	server.
		SetupMiddleware().
		FileServer().
		RegisterRoutes().
		Serve()
}
