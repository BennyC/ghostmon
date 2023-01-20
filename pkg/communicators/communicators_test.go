package communicators_test

import (
	"github.com/justpark/ghostmon/pkg/communicators"
	"github.com/justpark/ghostmon/pkg/logging"
	gtest "github.com/justpark/ghostmon/pkg/testing"
	"io"
	"net"
	"os"
	"testing"

	"github.com/sourcegraph/conc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommunicator_CutOver(t *testing.T) {
	server, client := net.Pipe()
	adapter := communicators.New(gtest.PipeConnector(client), logging.NewNilLogger())

	var wg conc.WaitGroup
	wg.Go(func() {
		err := adapter.Unpostpone()
		assert.NoError(t, err)
	})

	received, _ := io.ReadAll(server)
	require.Equal(t, "unpostpone", string(received))
	wg.Wait()
}

func TestCommunicator_CutOverErrors(t *testing.T) {
	server, client := net.Pipe()
	adapter := communicators.New(gtest.PipeConnector(client), logging.NewNilLogger())
	_ = server.Close()

	var wg conc.WaitGroup
	wg.Go(func() {
		err := adapter.Unpostpone()
		assert.Error(t, err)
	})

	_, _ = io.ReadAll(server)
	wg.Wait()
}

func TestCommunicator_Panic(t *testing.T) {
	server, client := net.Pipe()
	adapter := communicators.New(gtest.PipeConnector(client), logging.NewNilLogger())

	var wg conc.WaitGroup
	wg.Go(func() {
		err := adapter.Panic()
		assert.NoError(t, err)
	})

	received, _ := io.ReadAll(server)
	require.Equal(t, "panic", string(received))
	wg.Wait()
}

func TestCommunicator_Status(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	adapter := communicators.New(
		communicators.WithDialConnector(ln.Addr()),
		logging.NewNilLogger(),
	)

	var wg conc.WaitGroup
	wg.Go(func() {
		status, err := adapter.Status()

		require.NoError(t, err)
		require.NotNil(t, status)
		require.Contains(t, string(status.Body), "Copy: 0/2915")
	})

	// Server should receive the message "status" and reply with a multi-line response
	// regarding the current status of the gh-ost migration
	b := make([]byte, 6)
	w, _ := os.ReadFile("fixtures/status.txt")
	gtest.ListenAndRespond(t, ln, b, w)

	require.Equal(t, "status", string(b))
	wg.Wait()
}
