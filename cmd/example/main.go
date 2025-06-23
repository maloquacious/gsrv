package main

import (
	"log"
	"net/http"

	"github.com/maloquacious/gsrv"
)

func main() {
	// Create server with default settings
	server, err := gsrv.New(
		gsrv.WithHost("localhost"),
		gsrv.WithPort("3000"),
		gsrv.WithShutdownKey("andy"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Set up routes
	mux := http.NewServeMux()
	mux.Handle("GET /api/health", server.HealthHandler())
	mux.Handle("POST /api/shutdown/{key}", server.ShutdownHandler())
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	server.Handler = mux

	// Start server - blocks until shutdown signal received
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server error: %v", err)
	}
}
