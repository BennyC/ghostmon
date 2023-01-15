package ghostmon_test

import (
	"github.com/justpark/ghostmon"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLoadConfigDefaultValues(t *testing.T) {
	config, err := ghostmon.LoadConfig()

	require.NoError(t, err)
	require.Equal(t, config.ConnectionAddress, "localhost:9001")
	require.Equal(t, config.ConnectionType, "tcp")
}

func TestLoadConfigGrabsEnvsVars(t *testing.T) {
	_ = os.Setenv("CONNECTION_ADDR", "localhost:1001")
	_ = os.Setenv("CONNECTION_TYPE", "abc")
	config, err := ghostmon.LoadConfig()

	require.NoError(t, err)
	require.Equal(t, "localhost:1001", config.ConnectionAddress)
	require.Equal(t, "abc", config.ConnectionType)
}
