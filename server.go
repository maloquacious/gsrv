// Package gsrv implements a graceful shutdown wrapper around http.Server
package gsrv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ErrServerShutdown = Error("server shutdown")
)

// we assume that we're running behind a proxy for SSL

// New returns a new Server with default timeouts and graceful shutdown logic.
// Optional configuration can be applied via Option functions.
func New(options ...Option) (*Server, error) {
	s := &Server{}
	s.MaxHeaderBytes = 1 << 20
	s.IdleTimeout = 10 * time.Second
	s.ReadTimeout = 5 * time.Second
	s.WriteTimeout = 10 * time.Second

	// create a channel to listen for OS signals.
	s.admin.ctx = context.Background()
	s.admin.stop = make(chan os.Signal, 1)
	signal.Notify(s.admin.stop, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, err
		}
	}

	// create a random shutdown key if none was provided
	if s.admin.keys.shutdown == "" {
		s.admin.keys.shutdown = uuid.NewString()
	}

	return s, nil
}

type Server struct {
	http.Server
	host, port string
	admin      struct {
		ctx  context.Context
		stop chan os.Signal // channel to stop the server
		keys struct {
			shutdown string // key to stop the server
		}
	}
	started time.Time
}

func (s *Server) BaseURL() string {
	return fmt.Sprintf("%s://%s", "http", s.Addr)
}

// ListenAndServe starts the server and implements a graceful shutdown
func (s *Server) ListenAndServe() error {
	s.started = time.Now()
	// start the server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("listening on %s\n", s.BaseURL())
		// todo: should we return the error to the caller rather than exit?
		if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server: %v\n", err)
		}
		log.Printf("server: shutdown and closed\n")
	}()

	// server is running; block until we receive a signal.
	sig := <-s.admin.stop

	started := time.Now()
	log.Printf("signal: received %v (%v)\n", sig, time.Since(started))

	// graceful shutdown with a timeout.
	timeout := time.Second * 10
	log.Printf("creating context with %v timeout (%v)\n", timeout, time.Since(started))
	ctx, cancel := context.WithTimeout(s.admin.ctx, timeout)
	defer cancel()

	// cancel any idle connections.
	log.Printf("canceling idle connections (%v)\n", time.Since(started))
	s.SetKeepAlivesEnabled(false)

	log.Printf("sending signal to shut down the server (%v)\n", time.Since(started))
	if err := s.Shutdown(ctx); err != nil {
		return errors.Join(ErrServerShutdown, err)
	}

	log.Printf("server stopped ¡gracefully! (%v)\n", time.Since(started))
	return nil
}

// HealthHandler initiates a naive health check of the server.
// This expects that the route is defined as something like "GET /api/health."
func (s *Server) HealthHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		timeUp := time.Since(s.started)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]string{
			"uptime": timeUp.String(),
		}
		_ = json.NewEncoder(w).Encode(response)
	})
}

// ShutdownHandler initiates a graceful shutdown by sending SIGINT to the server.
// This expects that the route is defined as something like "POST /api/shutdown/{key}."
//
// todo: should throttle this
func (s *Server) ShutdownHandler() http.Handler {
	if s.ShutdownKey() == "" {
		log.Printf("[server] warning: installing handler before setting key")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else if s.ShutdownKey() == "" {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else if key := r.PathValue("key"); key != s.ShutdownKey() {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		response := map[string]string{
			"status": "server shutting down",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("[server] failed to write JSON response: %v", err)
		}

		// send shutdown signal to main server loop
		go s.TriggerShutdown()
	})
}

// ShutdownKey returns the key used by the graceful shutdown handler
func (s *Server) ShutdownKey() string {
	return s.admin.keys.shutdown
}

// TriggerShutdown safely signals the server to shut down.
// It will not block if the signal channel is already full.
func (s *Server) TriggerShutdown() {
	select {
	case s.admin.stop <- syscall.SIGINT:
		log.Printf("[server] shutdown signal sent")
	default:
		log.Printf("[server] shutdown signal already sent")
	}
}
