package ghostmon_test

import (
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
	adapter := ghostmon.NewCommunicator(client)

	var wg conc.WaitGroup
	wg.Go(func() {
		err := adapter.CutOver()
		assert.NoError(t, err)
	})

	received, _ := io.ReadAll(server)
	require.Equal(t, "unpostpone", string(received))
	wg.Wait()
}

func TestCommunicator_CutOverErrors(t *testing.T) {
	server, client := net.Pipe()
	adapter := ghostmon.NewCommunicator(client)
	_ = server.Close()

	var wg conc.WaitGroup
	wg.Go(func() {
		err := adapter.CutOver()
		assert.Error(t, err)
	})

	_, _ = io.ReadAll(server)
	wg.Wait()
}
