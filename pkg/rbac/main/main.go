package main

import (
	"github.com/patterninc/heimdall/pkg/rbac/ranger"
	"github.com/patterninc/heimdall/pkg/sql/parser/trino"
)

func main() {
	r := &ranger.ApacheRanger{
		Name:                  "LocalRanger",
		Client:                ranger.NewClient("http://localhost:6080", "admin", ""),
		ServiceName:           "TrinoRanger",
		AccessReceiver:        trino.NewTrinoAccessReceiver("glue_catalog"),
		SyncIntervalInMinutes: 5,
	}
	err := r.SyncState()
	if err != nil {
		panic(err)
	}

	println(r.HasAccess("ihladush", "SELECT * FROM test.test_table"))
	println(r.HasAccess("ivan.hladush@pattern.com", "SELECT * FROM test.test_table"))

}
