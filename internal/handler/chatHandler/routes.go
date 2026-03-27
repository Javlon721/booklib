package chatHandler

import (
	"fmt"
	"log/slog"

	chatbus "github.com/Javlon721/booklib/internal/bus/chatBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
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

	router.Post("/chat", authMid, handler.CreateChat)
	router.Post("/chat/:chatID", authMid, handler.CreateMessage)
	router.Get("/chat/:userID", authMid, handler.GetChatBy)
	router.Get("/messages/:chatID", authMid, handler.GetPendingMessages)

	// -------------------------------------------------------------------------
	router.Get("/ws/:chatID", authMid, func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("context", c.Context())
			return c.Next()
		}
		return errs.New(errs.UpgradeRequired, fmt.Errorf(""))
	}, websocket.New(handler.Websoket))

	router.Get("/stats/ws/:chatID", handler.WebsoketStats)
}
