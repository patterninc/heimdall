package server

import (
	"context"
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

func (s *Server) Start(ctx context.Context, router *mux.Router) error {

	srv := &http.Server{
		Addr:              s.Address,
		ReadTimeout:       time.Duration(s.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(s.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(s.IdleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(s.ReadHeaderTimeout) * time.Second,
		Handler:           router,
	}

	// run server in background
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// wait for either server error or context cancellation
	select {
	case err := <-srvErr:
		return err
	case <-ctx.Done():
		// give server some time to shutdown gracefully
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(timeoutCtx)
		return ctx.Err()
	}

}
