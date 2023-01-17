package testing

import (
	"github.com/justpark/ghostmon"
	"github.com/sourcegraph/conc"
	"io"
	"net"
	"net/http"
	"testing"
)

// ReceiveNetMessage will handle the next connection and receive the next message
// received on a net.Listener
func ReceiveNetMessage(t *testing.T, listener net.Listener) ([]byte, error) {
	t.Helper()

	var conn net.Conn
	var wg conc.WaitGroup

	wg.Go(func() {
		conn, _ = listener.Accept()
	})
	wg.Wait()

	b, _ := io.ReadAll(conn)
	return b, nil
}

// CreateHTTPServer will create a ghostmon.HTTPServer and give it the net address
// of a random testable TCP port. The listener for this port will also be
// returned to the caller
func CreateHTTPServer(t *testing.T) (*http.Server, net.Conn) {
	t.Helper()

	server, client := net.Pipe()
	communicator := ghostmon.NewCommunicator(PipeConnector(client))

	return ghostmon.NewHTTPServer(communicator), server
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
