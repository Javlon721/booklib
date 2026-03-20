package main

import (
	"log/slog"
	"os"
	"time"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	userStore "github.com/Javlon721/booklib/internal/bus/userBus/store"
	"github.com/Javlon721/booklib/internal/db/postgres"
	authhandler "github.com/Javlon721/booklib/internal/handler/authHandler"
	userhandler "github.com/Javlon721/booklib/internal/handler/userHandler"
	"github.com/Javlon721/booklib/internal/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	apiHost = ":8001"
)

func main() {
	// -------------------------------------------------------------------------
	// Initialize logger

	logger := slog.New(slog.Default().Handler())

	if err := run(logger); err != nil {
		logger.Error("startup", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	logger.Info("startup")

	// -------------------------------------------------------------------------
	// Initialize global config

	conf := struct {
		// Auth
		Token struct {
			Secret    string
			ExpiresAt time.Duration `envconfig:"EXPIRES_AT"`
		}

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

	// -------------------------------------------------------------------------
	// Initialize postgres connection

	logger.Info("startup", "status", "initizlizing db connection", "hostport", conf.Postgres.Host)

	dbConn, err := postgres.NewPool(postgres.Config(conf.Postgres))

	if err != nil {
		return err
	}

	// -------------------------------------------------------------------------
	// Initialize fiber app

	logger.Info("startup", "status", "initizlizing fiber app")

	var config = fiber.Config{ErrorHandler: middleware.Errors(logger)}

	app := fiber.New(config)

	// -------------------------------------------------------------------------
	// Add middlewares

	app.Use(middleware.Logger(logger))

	// -------------------------------------------------------------------------
	// Initialize API version: /api/v1

	appV1 := app.Group("/api/v1")

	// -------------------------------------------------------------------------
	// Initialize Handlers

	logger.Info("startup", "status", "initizlizing handlers")

	userStore := userStore.New(dbConn, logger)
	userBus := userbus.NewBusiness(logger, userStore)

	authhandler.Routes(appV1, userBus, logger, authhandler.TokenConfig{
		Secret:         []byte(conf.Token.Secret),
		TokenExpiresAt: conf.Token.ExpiresAt,
		Method:         jwt.SigningMethodHS256,
	})

	userhandler.Routes(appV1, userBus, logger)

	return app.Listen(apiHost)
}
