package heimdall

import (
	"context"
	"net/http"
	"time"

	"github.com/hladush/go-telemetry/pkg/telemetry"
)

type userNameType string

const (
	userNameKey userNameType = `username`
)

var (
	authTelemetryMethod = telemetry.NewMethod("auth", "Auth middleware")
)

func (h *Heimdall) auth(next http.Handler) http.Handler {
	// start latency timer
	defer authTelemetryMethod.RecordLatency(time.Now())

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// auth request count
		authTelemetryMethod.CountRequest()

		// let's get username from the header
		username := ``
		if h.Auth != nil {
			// TODO: process error here...
			username, _ = h.Auth.GetUser(r)
		}

		// let's write this username to request context...
		ctx := context.WithValue(r.Context(), userNameKey, username)

		// pass the chain...
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getUsername(r *http.Request) string {

	usernameValue := r.Context().Value(userNameKey)

	if usernameValue == nil {
		return ``
	}

	if username, ok := usernameValue.(string); ok {
		return username
	}

	return ``

}
