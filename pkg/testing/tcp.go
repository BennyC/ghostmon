package testing

import (
	"github.com/justpark/ghostmon/pkg/communicators"
	ghttp "github.com/justpark/ghostmon/pkg/http"
	"github.com/justpark/ghostmon/pkg/logging"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"testing"
)

var _ communicators.Connector = &pipeConnector{}

type pipeConnector struct {
	server net.Conn
}

// PipeConnector creates a net.Pipe and creates a communicators.Connector, allowing
// the two to communicate through our communicators.Communicator
func PipeConnector(s net.Conn) communicators.Connector {
	return &pipeConnector{server: s}
}

func (p pipeConnector) Connect() (net.Conn, error) {
	return p.server, nil
}

// CreateHTTPServer will create a ghostmon.HTTPServer and give it the net address
// of a random testable TCP port. The listener for this port will also be
// returned to the caller
func CreateHTTPServer(t *testing.T) (*http.Server, net.Conn) {
	t.Helper()
	logger := logging.NewNilLogger()
	server, client := net.Pipe()
	communicator := communicators.New(PipeConnector(client), logger)

	return ghttp.NewHTTPServer(":8080", communicator, logger), server
}

// ListenAndRespond will take a net.Listener and accept the next net.Conn. Will Read the contents of
// what was sent into the in parameter. If out is provided, it will Write the out parameter to the net.Conn
func ListenAndRespond(
	t *testing.T,
	ln net.Listener,
	in []byte,
	out ...[]byte,
) {
	t.Helper()

	conn, err := ln.Accept()
	require.NoError(t, err)

	_, err = conn.Read(in)
	require.NoError(t, err)

	if len(out) == 1 {
		_, err := conn.Write(out[0])
		require.NoError(t, err)
	}

	require.NoError(t, conn.Close())
}
