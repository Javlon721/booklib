package middleware

import (
	"fmt"
	"log/slog"

	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func WSUpgrade(loggers *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("context", c.Context())
			return c.Next()
		}
		return errs.New(errs.UpgradeRequired, fmt.Errorf(""))
	}
}
