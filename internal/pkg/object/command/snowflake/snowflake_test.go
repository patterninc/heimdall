package snowflake

import (
	"fmt"
	"strings"
	"testing"

	"github.com/patterninc/heimdall/pkg/plugin"
)

func TestResumable(t *testing.T) {

	h, err := New(nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := h.(plugin.ResultStreamer); !ok {
		t.Fatal("snowflake does not implement plugin.ResultStreamer")
	}
	if !plugin.IsResumable(h) {
		t.Fatal("snowflake is not resumable")
	}

}

func TestQueryTag(t *testing.T) {

	if got := queryTag(`abc-123`); got != `heimdall-abc-123` {
		t.Fatalf("got %q, want heimdall-abc-123", got)
	}

	query := fmt.Sprintf(queryHistorySelect, 1700000000-queryHistorySlack)
	if !strings.Contains(query, `to_timestamp_ltz(1699999940)`) || strings.Contains(query, `%!`) {
		t.Fatalf("query history select is not formatted as expected:\n%s", query)
	}

}
