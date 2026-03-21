package authhandler

import (
	"log/slog"

	authbus "github.com/Javlon721/booklib/internal/bus/authBus"
	"github.com/gofiber/fiber/v3"
)

func Routes(mux fiber.Router, userBus UserBus, logger *slog.Logger, auth *authbus.Bussiness) {
	router := mux.Group("/auth")

	handler := NewHandler(userBus, logger, auth)

	router.Post("/login", handler.Login)
}
