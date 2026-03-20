package middleware

import (
	"context"
	"log/slog"
	"strings"

	authhandler "github.com/Javlon721/booklib/internal/handler/authHandler"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/gofiber/fiber/v3"
)

type Authenticator interface {
	Authenticate(context.Context, string) (authhandler.AuthenticateResp, error)
}

func Authenticate(logger *slog.Logger, a Authenticator) fiber.Handler {
	return func(c fiber.Ctx) error {
		bearer := c.Get("Authorization")

		if !strings.HasPrefix(bearer, "Bearer ") {
			return errs.New(errs.Unauthenticated, authhandler.ErrInvalidCredentials)
		}

		tokenString := bearer[7:]

		ctx := c.Context()

		resp, err := a.Authenticate(ctx, tokenString)

		if err != nil {
			return err
		}

		ctx = setUserID(c.Context(), resp.UserID)

		c.SetContext(ctx)

		return c.Next()
	}
}
