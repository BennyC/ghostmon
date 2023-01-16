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

	comm, err := ghostmon.NewNetCommunicator(config)
	if err != nil {
		// TODO Handle err
	}

	server := ghostmon.NewHTTPServer(comm)
	if err := server.ListenAndServe(); err != nil {
		// TODO Handle err
	}
}
