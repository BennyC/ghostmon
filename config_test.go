package ghostmon_test

import (
	"github.com/justpark/ghostmon"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadConfigDefaultValues(t *testing.T) {
	config, err := ghostmon.LoadConfig()

	require.NoError(t, err)
	require.Equal(t, config.ConnectionAddress, "localhost:9001")
	require.Equal(t, config.ConnectionType, "tcp")
}
