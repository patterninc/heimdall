package main

import (
	"net/http"

	"github.com/patterninc/heimdall/pkg/plugin"
)

type authHeader struct {
	Header string `yaml:"header" json:"header"`
}

func New() (plugin.Auth, error) {
	return &authHeader{}, nil
}

func (a *authHeader) GetUser(r *http.Request) (string, error) {
	return r.Header.Get(a.Header), nil
}
