package main

import (
	"fmt"
	"monitoring/config"
	"monitoring/pkg/postgres"

	// "monitoring/internal/delivery/cron"
	rest "monitoring/internal/delivery/rest"
	. "monitoring/internal/globals"
	"monitoring/internal/util/midlog"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	GlobalConfig = config.NewConfig()

	SetGlobals()
	midlog.InfoF("Starting Midlog")
	midlog.LogCommandLine()
	if GlobalConfig.LogLevel != "" {
		midlog.InfoF("Setting log level to %v", GlobalConfig.LogLevel)
		logLevel, err := zerolog.ParseLevel(GlobalConfig.LogLevel)
		if err != nil {
			midlog.FatalF("Error parsing log level: %v", err)
		}
		midlog.SetLevel(logLevel)
	}

	// cronJob, err := cron.New()
	// if err != nil {
	// 	midlog.FatalF("Error creating cron job: %v", err)
	// }
	// cronJob.Start()

	r, err := rest.New()
	if err != nil {
		midlog.FatalF("Error creating rest server: %v", err)
	}
	err = r.Start(GlobalConfig.HTTP.Address)
	if err != nil {
		midlog.FatalF("Error starting rest server: %v", err)
	}

}

func SetGlobals() {

	var err error
	GlobalPG, err = postgres.New(GlobalConfig.PGConn)
	if err != nil {
		midlog.FatalF("Error creating midgard postgres client: %v", err)
	}

	GlobalConfig.HTTP.Address = "127.0.0.1:8090"
	GlobalConfig.HTTP.Debug = true

}
