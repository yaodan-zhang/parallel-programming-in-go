package main

import (
	"encoding/json"
	"os"
	"proj2/server"
	"strconv"
)

// The main function.
func main() {
	// Speficy input and output sources.
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	var mode string
	var numOfThreads int

	// Input didn't specify the number of threads, call sequential mode.
	if len(os.Args) == 1 {
		mode = "s"
	} else if len(os.Args) == 2 {
		// Input specified the number of threads, call parallel mode
		mode = "p"
		numOfThreads, _ = strconv.Atoi(os.Args[1])
	}

	var config = server.Config{Encoder: enc, Decoder: dec, Mode: mode, ConsumersCount: numOfThreads}

	server.Run(config)
}
