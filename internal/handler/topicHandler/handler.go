package topichandler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	topicbus "github.com/Javlon721/booklib/internal/bus/topicBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/handler/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler struct {
	topicBus *topicbus.Business
	logger   *slog.Logger
}

func NewHandler(topicBus *topicbus.Business, logger *slog.Logger) *Handler {
	return &Handler{
		topicBus: topicBus,
		logger:   logger,
	}
}

func (h Handler) Create(c fiber.Ctx) error {
	var nt NewTopic

	if err := c.Bind().Body(&nt); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	userID, err := middleware.GetUserID(c.Context())

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	busTopic, err := toBusTopic(nt, userID)

	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("cannot parse topic: %w", err))
	}

	topic, err := h.topicBus.Create(c.Context(), busTopic)

	if err != nil {
		if errors.Is(err, topicbus.ErrTopicAlreadyExists) {
			return errs.New(errs.AlreadyExists, err)
		}
		return errs.New(errs.Internal, err)
	}

	return c.Status(http.StatusCreated).JSON(toHandlerTopic(topic))
}

func (h Handler) GetTopicByID(c fiber.Ctx) error {
	topicID, err := uuid.Parse(c.Params("topicID"))

	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("cannot parse topicID: %w", err))
	}

	topic, err := h.topicBus.GetUserByID(c.Context(), topicID)

	if err != nil {
		if errors.Is(err, topicbus.ErrTopicNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return errs.New(errs.Internal, err)
	}

	return c.JSON(toHandlerTopic(topic))
}
