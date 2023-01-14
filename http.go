package ghostmon

// Endpoints
// - Status
// - CutOver
// - Abort

type HTTPServer struct {
	communication Communicator
}

func NewHTTPServer(adapter Communicator) *HTTPServer {
	return &HTTPServer{
		communication: adapter,
	}
}
