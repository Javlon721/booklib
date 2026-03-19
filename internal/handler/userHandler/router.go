package userhandler

import (
	"log/slog"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/gofiber/fiber/v3"
)

func Routes(mux fiber.Router, userBus *userbus.Business, logger *slog.Logger) {
	router := mux.Group("/users")

	userHandler := NewHandler(userBus, logger)

	router.Post("/", userHandler.Create)
	router.Get("/:userID", userHandler.GetUserByID)
}
