package config

import (
	"github.com/spf13/viper"
	"net"
)

var _ net.Addr = Config{}

type Config struct {
	ConnectionType    string `mapstructure:"CONNECTION_TYPE"`
	ConnectionAddress string `mapstructure:"CONNECTION_ADDR"`
	HTTPAddr          string `mapstructure:"HTTP_ADDR"`
}

func (c Config) Network() string {
	return c.ConnectionType
}

func (c Config) String() string {
	return c.ConnectionAddress
}

// Load will return a *Config with system required parameters
// Any failures to load the system configuration, an error will be
// returned to the caller
func Load() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetDefault("CONNECTION_TYPE", "tcp")
	viper.SetDefault("CONNECTION_ADDR", "localhost:9001")
	viper.SetDefault("HTTP_ADDR", ":8080")

	var c *Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return c, nil
}
