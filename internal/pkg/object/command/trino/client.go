package trino

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
	"github.com/patterninc/heimdall/pkg/result/column"
)

type request struct {
	endpoint  string
	headers   http.Header
	userAgent string
	client    *http.Client
	stdout    *os.File
	nextUri   string
	result    *result.Result
	state     string
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
	Columns []*trinoColumn `json:"columns,omitempty"`
	Data    [][]any        `json:"data,omitempty"`
	Stats   *stats         `json:"stats,omitempty"`
	Error   map[string]any `json:"error"`
}

func newRequest(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) (*request, error) {

	// get cluster context
	clusterCtx := &clusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterCtx); err != nil {
			return nil, err
		}
	}

	// get job context
	jobCtx := &jobContext{}
	if j.Context != nil {
		if err := j.Context.Unmarshal(jobCtx); err != nil {
			return nil, err
		}
	}
	jobCtx.Query = normalizeTrinoQuery(jobCtx.Query)

	// form context for trino request
	req := &request{
		endpoint: clusterCtx.Endpoint,
		headers: http.Header{
			`Content-Type`:    []string{`application/json`},
			`User-Agent`:      []string{r.UserAgent},
			`X-Trino-User`:    []string{j.User},
			`X-Trino-Source`:  []string{r.UserAgent},
			`X-Trino-Catalog`: []string{clusterCtx.Catalog},
		},
		userAgent: r.UserAgent,
		client:    &http.Client{},
		stdout:    r.Stdout,
		result:    &result.Result{},
	}

	// submit query
	if err := req.submit(jobCtx.Query); err != nil {
		return nil, err
	}

	return req, nil

}

func (r *request) submit(query string) error {

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/statement", r.endpoint), bytes.NewBuffer([]byte(query)))
	if err != nil {
		return err
	}

	return r.api(req)

}

func (r *request) poll() error {

	req, err := http.NewRequest(http.MethodGet, r.nextUri, nil)
	if err != nil {
		return err
	}

	return r.api(req)

}

func (r *request) api(req *http.Request) error {

	// let's set headers
	req.Header = r.headers

	// make api call
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code [%d]: %s", resp.StatusCode, string(body))
	}

	output := &response{}
	if err := json.Unmarshal(body, output); err != nil {
		return fmt.Errorf("error [%s]: %s", err, string(body))
	}

	// do we have an error?
	if output.Error != nil {
		errorMessage := ``
		if message, found := output.Error[`message`]; found {
			errorMessage = message.(string)
		} else {
			errorJson, err := json.Marshal(output.Error)
			if err != nil {
				return fmt.Errorf("cannot marshal: %v", output.Error)
			}
			errorMessage = string(errorJson)
		}
		return fmt.Errorf("query failed: %v", errorMessage)
	}

	// did we get columns metadata?
	if l := len(output.Columns); l > 0 && r.result.Columns == nil {
		r.result.Columns = make([]*column.Column, l)
		for i, c := range output.Columns {
			r.result.Columns[i] = &column.Column{Name: c.Name, Type: column.Type(c.Type)}
		}
	}

	// if we got any data, let's append it...
	if l := len(output.Data); l > 0 {
		if r.result.Data == nil {
			r.result.Data = make([][]any, 0, l*2) // we reserve "two pages of data" to potentially acount for additional rows...
		}
		r.result.Data = append(r.result.Data, output.Data...)
	}

	// do we have next page?
	r.nextUri = output.NextUri

	// let's record the last state we saw...
	if output.Stats != nil {
		r.state = output.Stats.State
	}

	return nil

}

func normalizeTrinoQuery(query string) string {
	// Trino does not support semicolon at the end of the query, so we remove it if present
	return strings.TrimSuffix(query, ";")
}
