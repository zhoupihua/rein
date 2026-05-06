package output

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	ExitOK         = 0
	ExitError      = 1
	ExitValidation = 2
	ExitNotFound   = 3
	ExitNoProject  = 4
)

func PrintJSON(data any) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling JSON: %v\n", err)
		os.Exit(ExitError)
	}
	fmt.Println(string(b))
}

func PrintHuman(data any) {
	// Default to JSON output — it's readable for both humans and AI
	PrintJSON(data)
}

func Print(data any, jsonMode bool) {
	if jsonMode {
		PrintJSON(data)
	} else {
		PrintHuman(data)
	}
}

func PrintError(err error, jsonMode bool) {
	if jsonMode {
		PrintJSON(map[string]string{"error": err.Error()})
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
	os.Exit(ExitError)
}
