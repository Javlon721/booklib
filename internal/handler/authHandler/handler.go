package authhandler

import (
	"context"
	"errors"
	"log/slog"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidAuthParams  = errors.New("invalid auth params")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserBus interface {
	GetUserByID(context.Context, uuid.UUID) (userbus.User, error)
	GetUserByEmail(context.Context, string) (userbus.User, error)
}

type Handler struct {
	userBus UserBus
	logger  *slog.Logger
}

func (h Handler) Login(c fiber.Ctx) error {
	var payload AuthParams

	if err := c.Bind().Body(&payload); err != nil {
		return errs.New(errs.InvalidArgument, ErrInvalidAuthParams)
	}

	user, err := h.userBus.GetUserByEmail(c.Context(), payload.Email)

	if err != nil {
		if errors.Is(err, userbus.ErrUserNotFound) {
			return errs.New(errs.NotFound, userbus.ErrUserNotFound)
		}

		return errs.New(errs.Internal, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		return errs.New(errs.PermissionDenied, ErrInvalidCredentials)
	}

	return c.Send([]byte("login succeded"))
}

func NewHandler(userBus UserBus, logger *slog.Logger) *Handler {
	return &Handler{
		userBus: userBus,
		logger:  logger,
	}
}
