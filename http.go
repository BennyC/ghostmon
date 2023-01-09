package ghostmon

// Endpoints
// - Status
// - CutOver
// - Abort

type HTTPServer struct {
	communication CommunicationAdapter
}

func NewHTTPServer(adapter CommunicationAdapter) *HTTPServer {
	return &HTTPServer{
		communication: adapter,
	}
}
