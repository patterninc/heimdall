package plugin

import (
	"net/http"
)

type Auth interface {
	GetUser(r *http.Request) (string, error)
}
