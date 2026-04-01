package authhandler

import (
	"context"
	"errors"
	"log/slog"

	authbus "github.com/Javlon721/booklib/internal/bus/authBus"
	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

var (
	ErrInvalidAuthParams = errors.New("invalid auth params")
)

type UserBus interface {
	GetUserByID(context.Context, uuid.UUID) (userbus.User, error)
	GetUserByEmail(context.Context, string) (userbus.User, error)
}

type Handler struct {
	userBus UserBus
	auth    *authbus.Bussiness
	logger  *slog.Logger
}

// Login
//
// @Summary Login user
// @Description Logins user
// @Tags auth
// @Param request body AuthParams true "User login request body"
// @Accept json
// @Produce plain
// @Success 200 {string} string
// @Failure 400 {object} errs.ErrResponce
// @Failure 401 {object} errs.ErrResponce
// @Failure 403 {object} errs.ErrResponce
// @Router /auth/login [post]
func (h Handler) Login(c fiber.Ctx) error {
	var payload AuthParams

	if err := c.Bind().Body(&payload); err != nil {
		return errs.New(errs.InvalidArgument, ErrInvalidAuthParams)
	}

	params, err := toBusAuthParams(payload)

	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	token, err := h.auth.Login(c.Context(), params)

	if err != nil {
		if errors.Is(err, authbus.ErrInvalidCredentials) {
			return errs.New(errs.PermissionDenied, err)
		}

		if errors.Is(err, userbus.ErrUserNotFound) {
			return errs.New(errs.Unauthenticated, err)
		}

		return errs.New(errs.Internal, err)
	}

	return c.SendString(token)
}

func NewHandler(userBus UserBus, logger *slog.Logger, auth *authbus.Bussiness) *Handler {
	return &Handler{
		userBus: userBus,
		logger:  logger,
		auth:    auth,
	}
}
