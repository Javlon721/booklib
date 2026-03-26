package discussionHandler

import (
	"log/slog"

	discussionsBus "github.com/Javlon721/booklib/internal/bus/discussionBus"
	"github.com/gofiber/fiber/v3"
)

func Routes(
	mux fiber.Router,
	discussionsBus *discussionsBus.Business,
	logger *slog.Logger,
	authMid fiber.Handler,
) {
	router := mux.Group("/topicMessages")

	handler := NewHandler(discussionsBus, logger)

	router.Post("/", authMid, handler.Create)
}
