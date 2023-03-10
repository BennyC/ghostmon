package main

import (
	"github.com/justpark/ghostmon/pkg/communicators"
	"github.com/justpark/ghostmon/pkg/config"
	"github.com/justpark/ghostmon/pkg/http"
	"github.com/justpark/ghostmon/pkg/logging"
	"golang.org/x/exp/slog"
)

func main() {
	logger := logging.New()
	logger.Info("loading configuration from environment vars")
	c, err := config.Load()
	if err != nil {
		logger.Error("unable to load configuration", err)
		return
	}

	comm := communicators.New(communicators.WithDialConnector(c), logger)
	server := http.NewHTTPServer(c.HTTPAddr, comm, logger)
	logger.Info("starting http server", slog.String("addr", c.HTTPAddr))
	if err := server.ListenAndServe(); err != nil {
		logger.Error("unable to start http server", err)
		return
	}
}
