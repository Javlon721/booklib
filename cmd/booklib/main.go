package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	authbus "github.com/Javlon721/booklib/internal/bus/authBus"
	chatbus "github.com/Javlon721/booklib/internal/bus/chatBus"
	chatStore "github.com/Javlon721/booklib/internal/bus/chatBus/store"
	discussionBus "github.com/Javlon721/booklib/internal/bus/discussionBus"
	discussionStore "github.com/Javlon721/booklib/internal/bus/discussionBus/store"
	topicbus "github.com/Javlon721/booklib/internal/bus/topicBus"
	topicStore "github.com/Javlon721/booklib/internal/bus/topicBus/store"
	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	userStore "github.com/Javlon721/booklib/internal/bus/userBus/store"
	"github.com/Javlon721/booklib/internal/db/postgres"
	authhandler "github.com/Javlon721/booklib/internal/handler/authHandler"
	"github.com/Javlon721/booklib/internal/handler/chatHandler"
	"github.com/Javlon721/booklib/internal/handler/chatWsHandler"
	discussionHandler "github.com/Javlon721/booklib/internal/handler/discussionHandler"
	"github.com/Javlon721/booklib/internal/handler/middleware"
	topichandler "github.com/Javlon721/booklib/internal/handler/topicHandler"
	userhandler "github.com/Javlon721/booklib/internal/handler/userHandler"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	_ "github.com/Javlon721/booklib/docs"
	swaggo "github.com/gofiber/contrib/v3/swaggo"
)

var (
	apiHost       = ":8001"
	shutdownTimer = time.Second * 15
)

// @title Booklib API
// @version 1.0
// @description This is a sample server
// @host 127.0.0.1:8001
// @BasePath /api/v1
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

	defer dbConn.Close()

	// -------------------------------------------------------------------------
	// Initialize fiber app

	logger.Info("startup", "status", "initizlizing fiber app")

	var config = fiber.Config{ErrorHandler: middleware.Errors(logger)}

	app := fiber.New(config)

	// -------------------------------------------------------------------------
	// Swagger
	app.Get("/swagger/*", swaggo.HandlerDefault)

	// -------------------------------------------------------------------------
	// Add middlewares

	app.Use(middleware.Logger(logger))

	// -------------------------------------------------------------------------
	// Initialize API version: /api/v1

	appV1 := app.Group("/api/v1")

	wd, _ := os.Getwd()
	staticDir := path.Join(wd, "static")

	app.Use("/static", static.New(staticDir))
	// -------------------------------------------------------------------------
	// Initialize Handlers

	logger.Info("startup", "status", "initizlizing handlers")

	userStore := userStore.New(dbConn, logger)
	userBus := userbus.NewBusiness(logger, userStore)

	auth := authbus.NewBussiness(logger, userBus, authbus.TokenConfig{
		Secret:         []byte(conf.Token.Secret),
		TokenExpiresAt: conf.Token.ExpiresAt,
		Method:         jwt.SigningMethodHS256,
	})

	topicStore := topicStore.New(dbConn, logger)
	topicBus := topicbus.NewBusiness(logger, topicStore)

	discussionStore := discussionStore.New(dbConn, logger)
	discussionBus := discussionBus.NewBusiness(logger, discussionStore)

	chatStore := chatStore.New(dbConn, logger)
	chatBus := chatbus.NewBusiness(logger, chatStore)

	// -------------------------------------------------------------------------
	// Auth Middleware

	authMid := middleware.Authenticate(logger, auth)
	authorizeMid := middleware.Authorize(logger, auth)

	// -------------------------------------------------------------------------
	// Register handlers routes

	authhandler.Routes(appV1, userBus, logger, auth)

	userhandler.Routes(appV1, userBus, logger, authMid, authorizeMid)

	topichandler.Routes(appV1, topicBus, logger, authMid)

	discussionHandler.Routes(appV1, discussionBus, logger, authMid)

	chatHandler.Routes(appV1, chatBus, logger, authMid)

	chatWsHandler.Routes(appV1, chatBus, logger, authMid)

	// -------------------------------------------------------------------------
	// Start API Service

	logger.Info("startup", "status", "initializing V1 API support")

	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info("startup", "status", "api router started", "host", apiHost)
		serverErrors <- app.Listen(apiHost)
	}()

	select {
	case err = <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-shutdown.Done():
		logger.Info("shutdown", "status", "statring")

		stop()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimer)
		defer cancel()

		if err := app.ShutdownWithContext(ctx); err != nil {
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

		logger.Info("shutdown", "status", "BYE!")
	}

	return nil
}
