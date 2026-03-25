package topicmessagehandler

import (
	"log/slog"
	"net/http"

	topicmessagebus "github.com/Javlon721/booklib/internal/bus/topicMessageBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/handler/middleware"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	topicMessageBus *topicmessagebus.Business
	logger          *slog.Logger
}

func NewHandler(topicMessageBus *topicmessagebus.Business, logger *slog.Logger) *Handler {
	return &Handler{
		topicMessageBus: topicMessageBus,
		logger:          logger,
	}
}

func (h Handler) Create(c fiber.Ctx) error {
	var ntm NewTopicMessage

	if err := c.Bind().Body(&ntm); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ctx := c.Context()

	userID, err := middleware.GetUserID(ctx)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	busTopicMessage, err := toBusTopicMessage(ntm, userID)

	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	topicMessage, err := h.topicMessageBus.Create(ctx, busTopicMessage)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return c.Status(http.StatusCreated).JSON(toHandlerTopicMessage(topicMessage))
}
