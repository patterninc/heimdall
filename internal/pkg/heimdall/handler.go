package heimdall

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/hladush/go-telemetry/pkg/telemetry"
)

const (
	contentTypeKey   = `Content-Type`
	contentTypeJSON  = `application/json`
	contentTypePlain = `text/plain`
	idKey            = `id`
	errorKey         = `error`
	jsonUserKey      = `user`
)

var (
	ErrNoCaller    = fmt.Errorf(`cannot identify caller -- access denied`)
	apiCallMetrics = telemetry.NewMethod("payload_handler", "API calls")
)

type hasID interface {
	GetID() string
}

func writeAPIError(w http.ResponseWriter, err error, obj any) {
	// API request error count
	apiCallMetrics.CountError()

	response := map[string]string{
		errorKey: err.Error(),
	}

	if objectWithID, ok := obj.(hasID); ok {
		response[idKey] = objectWithID.GetID()
	}

	responseJSON, _ := json.Marshal(response)

	w.Header().Add(contentTypeKey, contentTypeJSON)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(responseJSON)
}

func payloadHandler[T any](fn func(*T) (any, error)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// start latency timer
		defer apiCallMetrics.RecordLatency(time.Now())
		// let's read the request payload
		defer r.Body.Close()

		// API request count
		apiCallMetrics.CountRequest()

		data, err := io.ReadAll(r.Body)
		if err != nil {
			writeAPIError(w, err, nil)
			return
		}

		// let's unmarshal the payload
		var payload T
		if len(data) > 2 {
			if err := json.Unmarshal(data, &payload); err != nil {
				writeAPIError(w, err, nil)
				return
			}
		}

		// URL params and query override the JSON Payload
		if err := applyQuery(&payload, r); err != nil {
			writeAPIError(w, err, nil)
			return
		}

		// execute request
		result, err := fn(&payload)
		if err != nil {
			writeAPIError(w, err, result)
			return
		}

		// write result to the user
		resultJson, err := json.Marshal(result)
		if err != nil {
			writeAPIError(w, err, result)
			return
		}

		// return result
		w.Header().Add(contentTypeKey, contentTypeJSON)
		w.WriteHeader(http.StatusOK)
		w.Write(resultJson)

		// API request success count
		apiCallMetrics.CountSuccess()
	}

}

func applyQuery[T any](obj *T, r *http.Request) error {

	// let's get the vars
	vars := mux.Vars(r)

	// let's apply URL query
	for k, values := range r.URL.Query() {
		for _, v := range values {
			vars[k] = v
		}
	}

	// inject caller's username
	// we want this to be the last assignment to overwrite anything user passes
	username := getUsername(r)
	if username == `` && r.Method != http.MethodGet {
		return ErrNoCaller
	}
	vars[jsonUserKey] = username

	// now we'll marshal the map to json...
	data, err := json.Marshal(vars)
	if err != nil {
		return err
	}

	// and now we'll unmarshal the json we got into the struct we have
	return json.Unmarshal(data, obj)

}
