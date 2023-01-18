package communicators_test

import (
	"github.com/justpark/ghostmon/pkg/communicators"
	"github.com/justpark/ghostmon/pkg/logging"
	gtest "github.com/justpark/ghostmon/pkg/testing"
	"io"
	"net"
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
