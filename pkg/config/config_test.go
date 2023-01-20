package config_test

import (
	"github.com/justpark/ghostmon/pkg/config"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLoadConfigDefaultValues(t *testing.T) {
	c, err := config.Load()

	require.NoError(t, err)
	require.Equal(t, c.ConnectionAddress, "localhost:9001")
	require.Equal(t, c.ConnectionType, "tcp")
}

func TestLoadConfigGrabsEnvsVars(t *testing.T) {
	_ = os.Setenv("CONNECTION_ADDR", "localhost:1001")
	_ = os.Setenv("CONNECTION_TYPE", "abc")
	c, err := config.Load()

	require.NoError(t, err)
	require.Equal(t, "localhost:1001", c.ConnectionAddress)
	require.Equal(t, "abc", c.ConnectionType)
}
