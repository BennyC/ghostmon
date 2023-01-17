package ghostmon_test

import (
	gtest "github.com/justpark/ghostmon/pkg/testing"
	"io"
	"net"
	"testing"

	"github.com/justpark/ghostmon"
	"github.com/sourcegraph/conc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommunicator_CutOver(t *testing.T) {
	server, client := net.Pipe()
	adapter := ghostmon.NewCommunicator(gtest.PipeConnector(client))

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
	adapter := ghostmon.NewCommunicator(gtest.PipeConnector(client))
	_ = server.Close()

	var wg conc.WaitGroup
	wg.Go(func() {
		err := adapter.Unpostpone()
		assert.Error(t, err)
	})

	_, _ = io.ReadAll(server)
	wg.Wait()
}
