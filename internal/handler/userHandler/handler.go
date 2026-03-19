package userhandler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler struct {
	userBus *userbus.Business
	logger  *slog.Logger
}

func NewHandler(userBus *userbus.Business, logger *slog.Logger) *Handler {
	return &Handler{
		userBus: userBus,
		logger:  logger,
	}
}

func (h Handler) Create(c fiber.Ctx) error {
	var nu NewUser

	if err := c.Bind().Body(&nu); err != nil {
		return errs.New(errs.Internal, err)
	}

	bu, err := toBusUser(nu)

	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("cannot parse user: %w", err))
	}

	user, err := h.userBus.Create(c.Context(), bu)

	if err != nil {
		if errors.Is(err, userbus.ErrUserAlreadyExists) {
			return errs.New(errs.AlreadyExists, userbus.ErrUserAlreadyExists)
		}
		return errs.New(errs.Internal, err)
	}

	c.SendStatus(http.StatusCreated)

	return c.JSON(toHandlerUser(user))
}

func (h Handler) GetUserByID(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userID"))

	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("cannot parse userID: %w", err))
	}

	user, err := h.userBus.GetUserByID(c.Context(), userID)

	if err != nil {
		if errors.Is(err, userbus.ErrUserNotFound) {
			return errs.New(errs.NotFound, userbus.ErrUserNotFound)
		}
		return errs.New(errs.Internal, err)
	}

	return c.JSON(toHandlerUser(user))
}
