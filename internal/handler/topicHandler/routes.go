package topichandler

import (
	"log/slog"

	topicbus "github.com/Javlon721/booklib/internal/bus/topicBus"
	"github.com/gofiber/fiber/v3"
)

func Routes(
	mux fiber.Router,
	topicBus *topicbus.Business,
	logger *slog.Logger,
	authMid fiber.Handler,
) {
	router := mux.Group("/topics")

	topicHandler := NewHandler(topicBus, logger)

	router.Post("/", authMid, topicHandler.Create)
	router.Get("/", topicHandler.Query)
	router.Get("/:topicID", topicHandler.GetTopicByID)
}
