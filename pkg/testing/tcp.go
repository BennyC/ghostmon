package testing

import (
	"github.com/justpark/ghostmon"
	ghttp "github.com/justpark/ghostmon/pkg/http"
	"net"
	"net/http"
	"testing"
)

// CreateHTTPServer will create a ghostmon.HTTPServer and give it the net address
// of a random testable TCP port. The listener for this port will also be
// returned to the caller
func CreateHTTPServer(t *testing.T) (*http.Server, net.Conn) {
	t.Helper()

	server, client := net.Pipe()
	communicator := ghostmon.NewCommunicator(PipeConnector(client))

	return ghttp.NewHTTPServer(communicator), server
}

var _ ghostmon.Connector = &pipeConnector{}

type pipeConnector struct {
	server net.Conn
}

func (p pipeConnector) Connect() (net.Conn, error) {
	return p.server, nil
}

func PipeConnector(s net.Conn) ghostmon.Connector {
	return &pipeConnector{server: s}
}
