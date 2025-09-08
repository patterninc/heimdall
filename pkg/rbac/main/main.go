package main

import "github.com/patterninc/heimdall/pkg/rbac/ranger"

func main() {
	r := &ranger.ApacheRanger{
		URL:                   "http://localhost:6080",
		Username:              "admin",
		Password:              "",
		ServiceName:           "TrinoRanger",
		SyncIntervalInMinutes: 5,
	}
	err := r.SyncState()
	if err != nil {
		panic(err)
	}

}
