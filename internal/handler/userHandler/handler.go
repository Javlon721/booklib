package userhandler

import (
	"fmt"
	"log/slog"
	"net/http"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
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
		c.SendStatus(http.StatusInternalServerError)
		return err
	}

	bu, err := toBusUser(nu)

	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return err
	}

	user, err := h.userBus.Create(c.Context(), bu)

	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return err
	}

	c.SendStatus(http.StatusCreated)

	return c.JSON(toHandlerUser(user))
}

func (h Handler) GetUserByID(c fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("userID"))

	if err != nil {
		c.SendStatus(http.StatusBadRequest)
		return fmt.Errorf("cannot parse userID: %w", err)
	}

	user, err := h.userBus.GetUserByID(c.Context(), userID)

	if err != nil {
		c.SendStatus(http.StatusInternalServerError)
		return err
	}

	return c.JSON(toHandlerUser(user))
}
