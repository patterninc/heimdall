package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	Address           string `yaml:"addr,omitempty" json:"addr,omitempty"`
	ReadTimeout       int    `yaml:"read_timeout,omitempty" json:"read_timeout,omitempty"`
	WriteTimeout      int    `yaml:"write_timeout,omitempty" json:"write_timeout,omitempty"`
	IdleTimeout       int    `yaml:"idle_timeout,omitempty" json:"idle_timeout,omitempty"`
	ReadHeaderTimeout int    `yaml:"read_header_timeout,omitempty" json:"read_header_timeout,omitempty"`
}

// Start serves until ctx is done, then waits up to shutdownTimeout for in-flight requests.
func (s *Server) Start(ctx context.Context, router *mux.Router, shutdownTimeout time.Duration) error {

	// configure server
	server := &http.Server{
		Addr:              s.Address,
		ReadTimeout:       time.Duration(s.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(s.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(s.IdleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(s.ReadHeaderTimeout) * time.Second,
		Handler:           router,
	}

	shutdown := make(chan error, 1)
	context.AfterFunc(ctx, func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		shutdown <- server.Shutdown(ctx)
	})

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	if err := <-shutdown; err != nil {
		return fmt.Errorf("shutdown: http requests: %w", err)
	}

	return nil

}
