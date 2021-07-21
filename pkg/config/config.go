package config

import (
	"github.com/BurntSushi/toml"
)

// KafkaConfig represent config for Kafka Message Broker
type KafkaConfig struct {
	Host    string `toml:"host"`
	Topic   string `toml:"topic"`
	GroupID string `toml:"group_id"`
}

// Config represents mix of settings for the app.
type Config struct {
	Log         Log         `toml:"log"`
	KafkaConfig KafkaConfig `toml:"kafka"`
}

// Log describes settings for application logger.
type Log struct {
	Level string `toml:"level"`
	Mode  string `toml:"mode"`
}

// New creates reads application configuration from the file.
func New(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
