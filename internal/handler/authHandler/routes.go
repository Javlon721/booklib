package authhandler

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func Routes(mux fiber.Router, userBus UserBus, logger *slog.Logger) {
	router := mux.Group("/auth")

	handler := NewHandler(userBus, logger)

	router.Post("/login", handler.Login)
}
