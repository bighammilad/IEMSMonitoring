package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type (
	HTTPConf struct {
		Address string
		Debug   bool
	}

	Config struct {
		Debug    bool
		HTTP     HTTPConf
		PGConn   string
		Cron     CronConfig
		LogLevel string
	}

	CronConfig struct {
		CPUInterval string
	}
)

// NewConfig returns app config.
func NewConfig() (cfg Config) {
	// General
	err := envconfig.Process("monitoring", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	return
}
