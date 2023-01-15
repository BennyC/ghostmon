package ghostmon

import (
	"net/http"
	"time"
)

// Endpoints
// - Status
// - CutOver
// - Abort

func NewHTTPServer(adapter *Communicator) *http.Server {
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
func handleUnpostpone(adapter *Communicator) http.HandlerFunc {
	return onlyAllowMethod(http.MethodPost, func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Hello, World"))
	})
}

// onlyAllowMethod will perform checks on the http.Request Method and check if it
// equals the allowed method. When a disallowed method is presented, a 405 Status
// will be returned to the client
func onlyAllowMethod(allowedMethod string, next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != allowedMethod {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		next(writer, request)
	}
}
