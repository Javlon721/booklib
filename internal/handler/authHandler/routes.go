package authhandler

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func Routes(mux fiber.Router, userBus UserBus, logger *slog.Logger, cfg TokenConfig) {
	router := mux.Group("/auth")

	handler := NewHandler(userBus, logger, cfg)

	router.Post("/login", handler.Login)
}
