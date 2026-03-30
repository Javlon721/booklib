package chatWsHandler

import (
	"log/slog"

	chatbus "github.com/Javlon721/booklib/internal/bus/chatBus"
	"github.com/Javlon721/booklib/internal/handler/middleware"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func Routes(
	mux fiber.Router,
	chatBus *chatbus.Business,
	logger *slog.Logger,
	authMid fiber.Handler,
) {
	router := mux.Group("/ws/chats")

	handler := NewHandler(chatBus, logger)

	wsUpgradeMid := middleware.WSUpgrade(logger)

	router.Get("/:chatID", authMid, wsUpgradeMid, websocket.New(handler.Websoket))
	router.Get("/stats/:chatID", handler.WebsoketStats)
}
