package main

import (
	"github.com/justpark/ghostmon/pkg/communicators"
	"github.com/justpark/ghostmon/pkg/config"
	"github.com/justpark/ghostmon/pkg/http"
	"github.com/justpark/ghostmon/pkg/logging"
)

func main() {
	logger := logging.New()
	c, err := config.Load()
	if err != nil {
		logger.Error("unable to load configuration", err)
		return
	}

	comm := communicators.New(communicators.WithDialConnector(c), logger)
	server := http.NewHTTPServer(comm)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("unable to start http server", err)
		return
	}
}
