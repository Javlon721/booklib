package topicmessagehandler

import (
	"log/slog"

	topicmessagebus "github.com/Javlon721/booklib/internal/bus/topicMessageBus"
	"github.com/gofiber/fiber/v3"
)

func Routes(
	mux fiber.Router,
	topicMessageBus *topicmessagebus.Business,
	logger *slog.Logger,
	authMid fiber.Handler,
) {
	router := mux.Group("/topicMessages")

	handler := NewHandler(topicMessageBus, logger)

	router.Post("/", authMid, handler.Create)
}
