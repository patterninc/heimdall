package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/babourine/x/pkg/process"
	"github.com/patterninc/heimdall/pkg/result"
)

func main() {

	// Define a flag to specify the input type
	inputType := flag.String("type", "message", "Specify the input type: 'text' or 'json'")
	flag.Parse()

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		process.Bail(`stdio`, err)
	}

	var r *result.Result
	if *inputType == "json" {
		var jsonData map[string]any
		if err := json.Unmarshal(input, &jsonData); err != nil {
			process.Bail(`json`, fmt.Errorf("invalid JSON input: %w", err))
		}
		r, err = result.FromJson(jsonData)
	} else {
		r, err = result.FromMessage(string(input))
	}

	if err != nil {
		process.Bail(`result`, err)
	}

	data, err := json.Marshal(r)
	if err != nil {
		process.Bail(`json`, err)
	}

	if err := os.WriteFile(`result.json`, data, 0600); err != nil {
		process.Bail(`file`, err)
	}

}
