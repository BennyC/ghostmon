package http

import (
	"encoding/json"
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
	mux.HandleFunc("/status", handleStatus(adapter, logger))
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
	return OnlyAllowMethod(http.MethodPost, func(writer http.ResponseWriter, request *http.Request) {
		if err := adapter.Unpostpone(); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Info("handleUnpostpone: successfully sent unpostpone cmd to gh-ost")
		writer.WriteHeader(http.StatusCreated)
	})
}

// handleStatus creates a http.HandlerFunc that will accept GET requests and communicate with our gh-ost process
// to find out what our current status is. Status is then returned via JSON to the caller.
func handleStatus(adapter *communicators.Communicator, logger *slog.Logger) http.HandlerFunc {
	type statusResponse struct {
		FullStatus string `json:"full_status"`
		Table      string `json:"table,omitempty"`
	}

	return OnlyAllowMethod(http.MethodGet, func(writer http.ResponseWriter, request *http.Request) {
		status, err := adapter.Status()
		if err != nil {
			logger.Warn("handleStatus: unable to fetch status of the migration")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		table, _ := status.Table()
		b, err := json.Marshal(&statusResponse{
			FullStatus: string(status.Body),
			Table:      table,
		})

		if err != nil {
			logger.Error("handleStatus: failed to marshal response", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := writer.Write(b); err != nil {
			logger.Error("handleStatus: failed to write response", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	})
}

// handleAbort creates a http.HandlerFunc that will only accept POST requests from the client
// "panic" command will then be sent to the gh-ost process the *communicators.Communicator is
// working with
func handleAbort(adapter *communicators.Communicator, logger *slog.Logger) http.HandlerFunc {
	return OnlyAllowMethod(http.MethodPost, func(writer http.ResponseWriter, request *http.Request) {
		if err := adapter.Panic(); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
		}

		logger.Info("handleAbort: sent panic cmd to gh-ost, gh-ost will now abort migration")
		writer.WriteHeader(http.StatusCreated)
	})
}
