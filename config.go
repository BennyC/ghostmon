package ghostmon

import "github.com/spf13/viper"

type Config struct {
	ConnectionType    string `mapstructure:"CONNECTION_TYPE"`
	ConnectionAddress string `mapstructure:"CONNECTION_ADDR"`
}

// LoadConfig will return a *Config with system required parameters
// Any failures to load the system configuration, an error will be
// returned to the caller
func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetDefault("CONNECTION_TYPE", "tcp")
	viper.SetDefault("CONNECTION_ADDR", "localhost:9001")

	var c *Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}

	return c, nil
}
