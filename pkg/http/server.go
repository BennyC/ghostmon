package http

import (
	"github.com/justpark/ghostmon/pkg/communicators"
	"net/http"
	"time"
)

func NewHTTPServer(adapter *communicators.Communicator) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/unpostpone", handleUnpostpone(adapter))
	mux.HandleFunc("/status", handleUnpostpone(adapter))
	mux.HandleFunc("/abort", handleUnpostpone(adapter))

	return &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

// handleUnpostpone creates a http.HandlerFunc that will only accept POST requests
// and send an "unpostpone" command over the Communicator to the postponed gh-ost process
func handleUnpostpone(adapter *communicators.Communicator) http.HandlerFunc {
	return onlyAllowMethod(http.MethodPost, func(writer http.ResponseWriter, request *http.Request) {
		if err := adapter.Unpostpone(); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusCreated)
	})
}
