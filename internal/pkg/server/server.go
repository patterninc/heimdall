package server

import (
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

func (s *Server) Start(router *mux.Router) error {

	// configure server
	server := &http.Server{
		Addr:              s.Address,
		ReadTimeout:       time.Duration(s.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(s.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(s.IdleTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(s.ReadHeaderTimeout) * time.Second,
		Handler:           router,
	}

	return server.ListenAndServe()

}
