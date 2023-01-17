package main

import (
	"fmt"
	"github.com/justpark/ghostmon"
)

func main() {
	config, err := ghostmon.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("unable to load configuration: %w", err))
	}

	comm := ghostmon.NewCommunicator(ghostmon.WithDialConnector(config))

	server := ghostmon.NewHTTPServer(comm)
	if err := server.ListenAndServe(); err != nil {
		// TODO Handle err
	}
}
