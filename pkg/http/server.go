package http

import (
	"github.com/justpark/ghostmon/pkg/communicators"
	"golang.org/x/exp/slog"
	"net/http"
	"time"
)

func NewHTTPServer(
	addr string,
	adapter *communicators.Communicator,
	logger *slog.Logger,
) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/unpostpone", handleUnpostpone(adapter, logger))
	mux.HandleFunc("/status", handleStatus(adapter))
	mux.HandleFunc("/abort", handleAbort(adapter, logger))

	return &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

// handleUnpostpone creates a http.HandlerFunc that will only accept POST requests
// and send an "unpostpone" command over the Communicator to the postponed gh-ost process
func handleUnpostpone(adapter *communicators.Communicator, logger *slog.Logger) http.HandlerFunc {
	return onlyAllowMethod(http.MethodPost, func(writer http.ResponseWriter, request *http.Request) {
		logger.Info("sending unpostpone cmd to gh-ost")
		if err := adapter.Unpostpone(); err != nil {
			logger.Error("unable to send unpostpone cmd to connected gh-ost instance", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info("successfully sent unpostpone cmd to gh-ost")
		writer.WriteHeader(http.StatusCreated)
	})
}

func handleStatus(adapter *communicators.Communicator) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNotImplemented)
	}
}

// handleAbort creates a http.HandlerFunc that will only accept POST requests from the client
// "panic" command will then be sent to the gh-ost process the *communicators.Communicator is
// working with
func handleAbort(adapter *communicators.Communicator, logger *slog.Logger) http.HandlerFunc {
	return onlyAllowMethod(http.MethodPost, func(writer http.ResponseWriter, request *http.Request) {
		logger.Info("sending panic cmd to gh-ost")
		if err := adapter.Panic(); err != nil {
			logger.Error("failed to send panic cmd to gh-ost", err)
			writer.WriteHeader(http.StatusInternalServerError)
		}

		logger.Info("sent panic cmd to gh-ost, gh-ost will now abort migration")
		writer.WriteHeader(http.StatusCreated)
	})
}
