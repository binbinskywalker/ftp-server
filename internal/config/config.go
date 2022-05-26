package config

// Version defines the ChirpStack Network Server version.
var Version string

// Config defines the configuration structure.
type Config struct {
	Ftp struct {
		User     string `mapstructure:"user"`
		PassWord string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Group    string `mapstructure:"group"`
		Owner    string `mapstructure:"owner"`
		DataPort string `mapstructure:"passive-port"`
	} `mapstructure:"ftp"`
}

// C holds the global configuration.
var C Config

// Get returns the configuration.
func Get() *Config {
	return &C
}

// Set sets the configuration.
func Set(c Config) {
	C = c
}
