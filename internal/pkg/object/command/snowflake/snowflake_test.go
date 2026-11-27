package snowflake

import (
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
	if _, ok := h.(plugin.Resumer); !ok {
		t.Fatal("snowflake does not implement plugin.Resumer")
	}

}

func TestQueryTag(t *testing.T) {

	if got := queryTag(`abc-123`); got != `heimdall-job-abc-123` {
		t.Fatalf("got %q, want heimdall-job-abc-123", got)
	}

}
