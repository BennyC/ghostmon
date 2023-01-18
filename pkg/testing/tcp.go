package testing

import (
	"github.com/justpark/ghostmon/pkg/communicators"
	ghttp "github.com/justpark/ghostmon/pkg/http"
	"github.com/justpark/ghostmon/pkg/logging"
	"net"
	"net/http"
	"testing"
)

// CreateHTTPServer will create a ghostmon.HTTPServer and give it the net address
// of a random testable TCP port. The listener for this port will also be
// returned to the caller
func CreateHTTPServer(t *testing.T) (*http.Server, net.Conn) {
	t.Helper()
	logger := logging.NewNilLogger()
	server, client := net.Pipe()
	communicator := communicators.New(PipeConnector(client), logger)

	return ghttp.NewHTTPServer(communicator, logger), server
}

var _ communicators.Connector = &pipeConnector{}

type pipeConnector struct {
	server net.Conn
}

func (p pipeConnector) Connect() (net.Conn, error) {
	return p.server, nil
}

func PipeConnector(s net.Conn) communicators.Connector {
	return &pipeConnector{server: s}
}
