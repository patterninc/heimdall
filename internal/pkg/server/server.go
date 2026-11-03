package server

import (
	"context"
	"errors"
	"net"
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
	server            *http.Server
	listener          net.Listener
}

// Listen binds the address. It must be called before Serve and Shutdown.
func (s *Server) Listen(router *mux.Router) error {

	// configure server
	s.server = &http.Server{
		Addr:              s.Address,
		ReadTimeout:       time.Duration(s.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(s.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(s.IdleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(s.ReadHeaderTimeout) * time.Second,
		Handler:           router,
	}

	listener, err := net.Listen(`tcp`, s.Address)
	if err != nil {
		return err
	}
	s.listener = listener

	return nil

}

// Serve blocks until the server fails or Shutdown is called.
func (s *Server) Serve() error {

	if err := s.server.Serve(s.listener); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil

}

// Shutdown stops accepting connections and waits for in-flight requests until ctx is done.
func (s *Server) Shutdown(ctx context.Context) error {

	return s.server.Shutdown(ctx)

}
