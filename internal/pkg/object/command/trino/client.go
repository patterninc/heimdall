package trino

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	failedState string = "FAILED"

	pollInterval = time.Second
)

type trinoClient struct {
	endpoint string
	user     string
	token    *string

	client *http.Client

	stdout *os.File
}

type trinoColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type stats struct {
	// Only care about State for now
	State string `json:"State"`
}

type response struct {
	NextUri string         `json:"nextUri,omitempty"`
	Columns []trinoColumn  `json:"columns,omitempty"`
	Data    [][]any        `json:"data,omitempty"`
	Stats   stats          `json:"stats,omitempty"`
	Error   map[string]any `json:"error"`
}

func newTrinoClient(endpoint, user string, token *string, r *plugin.Runtime) *trinoClient {
	return &trinoClient{
		endpoint: endpoint,
		user:     user,
		token:    token,
		client:   http.DefaultClient,
		stdout:   r.Stdout,
	}
}

func (t *trinoClient) Query(query string) (*response, error) {
	fmt.Fprintf(t.stdout, "Query: %s\n", query)
	finalOutput := &response{
		Columns: []trinoColumn{},
		Data:    [][]any{},
	}

	resp, err := t.start(query)
	if err != nil {
		return nil, err
	}

	finalOutput.Columns = resp.Columns
	finalOutput.Data = append(finalOutput.Data, resp.Data...)

	for {
		if resp.NextUri == "" {
			break
		}
		resp, err = t.poll(resp.NextUri)
		if err != nil {
			return nil, err
		}

		if resp.Stats.State == failedState {
			errMessage, _ := json.Marshal(resp.Error)
			return nil, errors.New(string(errMessage))
		}

		finalOutput.Columns = resp.Columns
		finalOutput.Data = append(finalOutput.Data, resp.Data...)

		time.Sleep(pollInterval)
	}

	return finalOutput, nil
}

func (t *trinoClient) start(query string) (*response, error) {
	t.stdout.WriteString("Starting trino query\n")
	req, err := t.newStartRequest(query)
	if err != nil {
		return nil, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return newResponse(resp)
}

func (t *trinoClient) poll(uri string) (*response, error) {
	fmt.Fprintf(t.stdout, "Polling trino query: %s\n", uri)
	req, err := t.newPollRequest(uri)
	if err != nil {
		return nil, err
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return newResponse(resp)
}

func (t *trinoClient) newStartRequest(query string) (*http.Request, error) {
	buffer := bytes.NewBuffer([]byte(query))
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/statement", t.endpoint), buffer)
	if err != nil {
		return nil, err
	}
	addHeaders(req, t.user, t.token)
	return req, nil
}

func (t *trinoClient) newPollRequest(nextUri string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, nextUri, nil)
	if err != nil {
		return nil, err
	}
	addHeaders(req, t.user, t.token)
	return req, nil
}

func addHeaders(req *http.Request, user string, token *string) {
	header := http.Header{
		"X-Presto-User": []string{user},
		"X-Trino-User":  []string{user},
	}
	if token != nil {
		header.Add("Authorization", *token)
	}

	req.Header = header
}

func newResponse(resp *http.Response) (*response, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	output := &response{}
	if err := json.Unmarshal(body, output); err != nil {
		return nil, err
	}

	return output, nil
}
