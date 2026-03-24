package middleware

import (
	"context"
	"fmt"
	"log/slog"

	authbus "github.com/Javlon721/booklib/internal/bus/authBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/types/role"
	"github.com/gofiber/fiber/v3"
)

type Authenticator interface {
	Authenticate(context.Context, string) (authbus.AuthenticateResp, error)
	Authorize(context.Context, []role.Role, []role.Role) bool
}

func Authenticate(logger *slog.Logger, a Authenticator) fiber.Handler {
	return func(c fiber.Ctx) error {
		bearerToken := c.Get("Authorization")

		if bearerToken == "" {
			return errs.New(errs.Unauthenticated, fmt.Errorf("authorization token not found"))
		}

		ctx := c.Context()

		resp, err := a.Authenticate(ctx, bearerToken)

		if err != nil {
			return errs.New(errs.Unauthenticated, err)
		}

		ctx = setUserID(c.Context(), resp.UserID)
		ctx = setUserRoles(ctx, resp.Roles)

		c.SetContext(ctx)

		return c.Next()
	}
}

type AuthorizeHandler func(roles ...role.Role) fiber.Handler

func Authorize(logger *slog.Logger, a Authenticator) AuthorizeHandler {
	w := func(roles ...role.Role) fiber.Handler {
		return func(c fiber.Ctx) error {
			userRoles, err := GetUserRoles(c.Context())

			if err != nil {
				return errs.New(errs.PermissionDenied, err)
			}

			isAuthorized := a.Authorize(c.Context(), roles, userRoles)

			if !isAuthorized {
				return errs.New(errs.PermissionDenied, fmt.Errorf("not enough roles"))
			}

			return c.Next()
		}
	}

	return w
}
