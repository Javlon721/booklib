package main

import (
	"log/slog"
	"os"
	"time"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	userStore "github.com/Javlon721/booklib/internal/bus/userBus/store"
	"github.com/Javlon721/booklib/internal/db/postgres"
	userhandler "github.com/Javlon721/booklib/internal/handler/userHandler"
	"github.com/Javlon721/booklib/internal/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	apiHost = ":8001"
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

	logger.Info("startup", "status", "initizlizing handlers")

	userStore := userStore.New(dbConn, logger)
	userBus := userbus.NewBusiness(logger, userStore)
	userHandler := userhandler.NewHandler(userBus, logger)

	var config = fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			if c.Response().StatusCode() >= 500 {
				return c.JSON(map[string]string{"error": "internal server error"})
			}

			return c.JSON(map[string]string{"error": err.Error()})
		},
	}

	logger.Info("startup", "status", "initizlizing fiber app")

	app := fiber.New(config)
	app.Use(middleware.Logger(logger))

	appV1 := app.Group("/api/v1")

	appV1.Post("/users", userHandler.Create)

	return app.Listen(apiHost)
}
