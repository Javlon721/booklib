package chatHandler

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
	router := mux.Group("/chats")

	handler := NewHandler(chatBus, logger)

	wsUpgradeMid := middleware.WSUpgrade(logger)

	router.Post("/chat", authMid, handler.CreateChat)
	router.Post("/chat/:chatID", authMid, handler.CreateMessage)
	router.Get("/chat/:userID", authMid, handler.GetChatBy)
	router.Get("/messages/:chatID", authMid, handler.GetPendingMessages)

	// -------------------------------------------------------------------------
	router.Get("/ws/:chatID", authMid, wsUpgradeMid, websocket.New(handler.Websoket))
	router.Get("/stats/ws/:chatID", handler.WebsoketStats)
}
