package userhandler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/handler/middleware"
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

// Create
//
// @Summary Create user
// @Description Creates new user
// @Tags users
// @Param request body NewUser true "Create user request body"
// @Accept json
// @Produce json
// @Success 201 {object} User
// @Failure 400 {object} errs.ErrResponce
// @Failure 409 {object} errs.ErrResponce
// @Router /users [post]
func (h Handler) Create(c fiber.Ctx) error {
	var nu NewUser

	if err := c.Bind().Body(&nu); err != nil {
		return errs.New(errs.InvalidArgument, err)
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

// GetUserByID
//
// @Summary Get user by id
// @Description Gets user by id
// @Tags users
// @Param Authorization header string true "JWT token"
// @Param userID path string true "User uuid"
// @Produce json
// @Success 200 {object} User
// @Failure 400 {object} errs.ErrResponce
// @Failure 401 {object} errs.ErrResponce
// @Failure 403 {object} errs.ErrResponce
// @Failure 404 {object} errs.ErrResponce
// @Router /users/{userID} [get]
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

// GetMe
//
// @Summary Get current user
// @Description Gets user from "subject" field in JWT token
// @Tags users
// @Param Authorization header string true "JWT token"
// @Produce json
// @Success 201 {object} User
// @Failure 401 {object} errs.ErrResponce
// @Failure 403 {object} errs.ErrResponce
// @Failure 404 {object} errs.ErrResponce
// @Router /users/me [get]
func (h Handler) GetMe(c fiber.Ctx) error {
	userID, err := middleware.GetUserID(c.Context())

	if err != nil {
		return err
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
