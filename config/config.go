package config

import (
	"encoding/json"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var cfg *Configuration

// Configuration structure which hold information for configuring the skills matrix API
type Configuration struct {
	BindAddr                string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	MongoConfig             MongoConfig
}

// MongoConfig contains the config required to connect to MongoDB.
type MongoConfig struct {
	BindAddr   string `envconfig:"MONGODB_BIND_ADDR"   json:"-"`
	Collection string `envconfig:"MONGODB_COLLECTION"`
	Database   string `envconfig:"MONGODB_DATABASE"`
}

// Get the application and returns the configuration structure
func Get() (*Configuration, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Configuration{
		BindAddr:                ":10000",
		GracefulShutdownTimeout: 5 * time.Second,
		MongoConfig: MongoConfig{
			BindAddr:   "localhost:27017",
			Collection: "timesheets",
			Database:   "timesheets",
		},
	}

	return cfg, envconfig.Process("", cfg)
}

// String is implemented to prevent sensitive fields being logged.
// The config is returned as JSON with sensitive fields omitted.
func (config Configuration) String() string {
	json, _ := json.Marshal(config)
	return string(json)
}
