package heimdall

import (
	"context"
	"net/http"
)

type userNameType string

const (
	userNameKey userNameType = `username`
)

func (h *Heimdall) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
