package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
)

func Logger(logger *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		err := c.Next()

		var (
			start              = time.Now()
			method             = c.Method()
			path               = c.Path()
			responceStatusCode = c.Response().StatusCode()
			ip                 = c.IP()
		)

		logger.Info(fmt.Sprintf("request: %s %s %s %d %v", ip, method, path, responceStatusCode, time.Since(start)))

		return err
	}
}
