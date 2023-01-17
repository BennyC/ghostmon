package main

import (
	"fmt"
	"github.com/justpark/ghostmon"
	"github.com/justpark/ghostmon/pkg/communicators"
	"github.com/justpark/ghostmon/pkg/http"
)

func main() {
	config, err := ghostmon.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("unable to load configuration: %w", err))
	}

	comm := communicators.New(communicators.WithDialConnector(config))

	server := http.NewHTTPServer(comm)
	if err := server.ListenAndServe(); err != nil {
		// TODO Handle err
	}
}
