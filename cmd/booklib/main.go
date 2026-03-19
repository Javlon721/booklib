package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/Javlon721/booklib/internal/db/postgres"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	logger := slog.New(slog.Default().Handler())

	if err := run(logger); err != nil {
		logger.Error("startup", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	logger.Info("startup")

	conf := struct {
		//Database
		Postgres struct {
			DB_Name  string
			Password string
			User     string
			Host     string

			MaxConns        int32         `envconfig:"MAX_CONNS"`
			MinConns        int32         `envconfig:"MIN_CONNS"`
			MaxConnIdleTime time.Duration `envconfig:"MAX_CONN_IDLE_TIME"`
			MaxConnLifetime time.Duration `envconfig:"MAX_CONN_LIFETIME"`
		}
	}{}

	err := godotenv.Load(".env")

	if err != nil {
		return err
	}

	err = envconfig.Process("", &conf)

	if err != nil {
		return err
	}

	logger.Info("startup", "status", "initizlizing db connection", "hostport", conf.Postgres.Host)

	dbConn, err := postgres.NewPool(postgres.Config(conf.Postgres))

	if err != nil {
		return err
	}

	_ = dbConn

	return nil
}
