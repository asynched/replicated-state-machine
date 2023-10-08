package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Host      string       `toml:"host"`
	Port      int          `toml:"port"`
	Bootstrap bool         `toml:"bootstrap"`
	Peers     []NodeConfig `toml:"peers"`
}

type NodeConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func ParseConfig(filename string) (Config, error) {
	config := Config{}

	if _, err := toml.DecodeFile(filename, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
