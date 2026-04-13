package userhandler

import (
	"log/slog"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/Javlon721/booklib/internal/handler/middleware"
	"github.com/Javlon721/booklib/internal/types/role"
	"github.com/gofiber/fiber/v3"
)

func Routes(
	mux fiber.Router,
	userBus *userbus.Business,
	logger *slog.Logger,
	authMid fiber.Handler,
	authorizeHandler middleware.AuthorizeHandler,
) {
	router := mux.Group("/users")

	authorizeAdmin := authorizeHandler(role.Admin)

	userHandler := NewHandler(userBus, logger)

	router.Post("/", userHandler.Create)
	router.Get("/me", authMid, userHandler.GetMe)
	router.Get("/:userID", authMid, authorizeAdmin, userHandler.GetUserByID)
	router.Put("/", authMid, userHandler.Update)
}
