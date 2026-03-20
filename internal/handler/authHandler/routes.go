package authhandler

import (
	"log/slog"

	authbus "github.com/Javlon721/booklib/internal/bus/authBus"
	"github.com/gofiber/fiber/v3"
)

func Routes(mux fiber.Router, userBus UserBus, logger *slog.Logger, cfg TokenConfig, auth *authbus.Bussiness) *Handler {
	router := mux.Group("/auth")

	handler := NewHandler(userBus, logger, cfg, auth)

	router.Post("/login", handler.Login)

	return handler
}
