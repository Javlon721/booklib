package userhandler

import (
	"log/slog"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/gofiber/fiber/v3"
)

func Routes(mux fiber.Router, userBus *userbus.Business, logger *slog.Logger, authMid fiber.Handler) {
	router := mux.Group("/users")

	userHandler := NewHandler(userBus, logger)

	router.Post("/", userHandler.Create)
	router.Get("/me", authMid, userHandler.GetMe)
	router.Get("/:userID", userHandler.GetUserByID)
}
