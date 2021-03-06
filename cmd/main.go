package main

import (
	"github.com/volodimyr/publisher/pkg/api/listener"
	"github.com/volodimyr/publisher/pkg/api/publisher"
	"github.com/volodimyr/publisher/pkg/persistence"
	"github.com/volodimyr/publisher/pkg/server"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "server: ", log.LstdFlags|log.Lshortfile)
	mux := http.NewServeMux()

	storage := persistence.New(logger)
	listener.NewHandlers(logger, storage).SetupRoutes(mux)
	publisher.NewHandlers(logger, storage).SetupRoutes(mux)

	ser := server.New(mux, ":8080")
	logger.Printf("Starting server at [%v] \n", ":8080")
	if err := ser.ListenAndServe(); err != nil {
		logger.Fatalf("server: failed to start [%v]\n", err)
	}
}
