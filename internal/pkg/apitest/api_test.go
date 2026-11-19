// Package apitest drives a running Heimdall server over HTTP, the way a client does.
// The tests need the server URL in HEIMDALL_TEST_URL and the ping command and localhost cluster from configs/local.yaml.
package apitest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	testServerEnvVar = `HEIMDALL_TEST_URL`
	userHeader       = `X-Heimdall-User`
	testUser         = `apitest_user`
	otherUser        = `apitest_other_user`
	apiPrefix        = `/api/v1`
	jobTimeout       = 60 * time.Second
	pollInterval     = 500 * time.Millisecond
	statusSucceeded  = `SUCCEEDED`
)

type jobResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

// TestPingJobLifecycle submits a ping job, waits for it to succeed and reads its result,
// which covers the API, the database, the async pool, plugin loading and result storage.
func TestPingJobLifecycle(t *testing.T) {

	baseURL := serverURL(t)

	code, body := call(t, http.MethodPost, baseURL+`/job`, testUser, `{
		"name": "apitest-ping",
		"version": "0.0.1",
		"context": {},
		"command_criteria": ["type:ping"],
		"cluster_criteria": ["type:localhost"]
	}`)
	if code != http.StatusOK {
		t.Fatalf("submit job: status %d, body %s", code, body)
	}

	submitted := parseJob(t, body)
	if submitted.ID == `` {
		t.Fatalf("submit job: no job id in %s", body)
	}

	// the job is async unless the command says otherwise, so wait for the pool to run it
	deadline := time.Now().Add(jobTimeout)
	for {
		code, body = call(t, http.MethodGet, baseURL+`/job/`+submitted.ID, testUser, ``)
		if code != http.StatusOK {
			t.Fatalf("get job: status %d, body %s", code, body)
		}
		j := parseJob(t, body)
		if j.Status == statusSucceeded {
			break
		}
		if j.Status == `FAILED` || j.Status == `KILLED` || j.Status == `CANCELED` {
			t.Fatalf("job ended as %s: %s", j.Status, j.Error)
		}
		if time.Now().After(deadline) {
			t.Fatalf("job still %s after %v", j.Status, jobTimeout)
		}
		time.Sleep(pollInterval)
	}

	resultURL := baseURL + `/job/` + submitted.ID + `/result`

	code, body = call(t, http.MethodGet, resultURL, testUser, ``)
	if code != http.StatusOK {
		t.Fatalf("get result: status %d, body %s", code, body)
	}
	if want := `Hello, ` + testUser + `!`; !strings.Contains(string(body), want) {
		t.Errorf("get result: body %s does not contain %q", body, want)
	}

	// only the user who submitted the job may read its result
	if code, body = call(t, http.MethodGet, resultURL, otherUser, ``); code != http.StatusForbidden {
		t.Errorf("get result as another user: status %d, want %d, body %s", code, http.StatusForbidden, body)
	}

}

// TestJobWithoutMatchingCommandIsRejected checks that criteria no command satisfies fail at submit.
func TestJobWithoutMatchingCommandIsRejected(t *testing.T) {

	baseURL := serverURL(t)

	code, body := call(t, http.MethodPost, baseURL+`/job`, testUser, `{
		"name": "apitest-no-match",
		"version": "0.0.1",
		"context": {},
		"command_criteria": ["type:apitest-does-not-exist"],
		"cluster_criteria": ["type:localhost"]
	}`)
	if code == http.StatusOK {
		t.Fatalf("submit job: status %d, want an error, body %s", code, body)
	}
	if want := `command-cluster pair is not found`; !strings.Contains(string(body), want) {
		t.Errorf("submit job: body %s does not contain %q", body, want)
	}

}

// TestClustersHealth checks that every cluster opted into health checks reports healthy.
func TestClustersHealth(t *testing.T) {

	baseURL := serverURL(t)

	if code, body := call(t, http.MethodGet, baseURL+`/clusters/health`, testUser, ``); code != http.StatusOK {
		t.Fatalf("clusters health: status %d, body %s", code, body)
	}

}

func serverURL(t *testing.T) string {

	t.Helper()

	url := os.Getenv(testServerEnvVar)
	if url == `` {
		t.Skipf("%s is not set", testServerEnvVar)
	}

	return strings.TrimRight(url, `/`) + apiPrefix

}

func call(t *testing.T, method, url, user, payload string) (int, []byte) {

	t.Helper()

	req, err := http.NewRequest(method, url, bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(userHeader, user)
	req.Header.Set(`Content-Type`, `application/json`)

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	return resp.StatusCode, body

}

func parseJob(t *testing.T, body []byte) *jobResponse {

	t.Helper()

	j := &jobResponse{}
	if err := json.Unmarshal(body, j); err != nil {
		t.Fatalf("unexpected response %s: %v", body, err)
	}

	return j

}
