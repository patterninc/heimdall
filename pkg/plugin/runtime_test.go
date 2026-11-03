package plugin

import (
	"encoding/json"
	"testing"
)

// the shell plugin writes the runtime into the job's context.json
func TestRuntimeMarshalsToJSON(t *testing.T) {

	r := &Runtime{WorkingDirectory: `/tmp/job`, Background: func(f func()) { go f() }}
	if _, err := json.Marshal(r); err != nil {
		t.Fatalf("marshal runtime: %v", err)
	}

}
