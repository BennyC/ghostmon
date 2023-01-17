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

// CreateServer will create a ghostmon.HTTPServer and give it the net address
// of a random testable TCP port. The listener for this port will also be
// returned to the caller
func CreateServer(t *testing.T) (*http.Server, net.Listener) {
	t.Helper()

	// TODO Handle err
	ln, _ := net.Listen("tcp", "127.0.0.1:0")

	communicator, err := ghostmon.NewNetCommunicator(ln.Addr())
	if err != nil {
		t.FailNow()
	}

	return ghostmon.NewHTTPServer(communicator), ln
}
