package discussionHandler

import (
	"log/slog"
	"net/http"

	discussionsBus "github.com/Javlon721/booklib/internal/bus/discussionBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/handler/middleware"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	discussionsBus *discussionsBus.Business
	logger         *slog.Logger
}

func NewHandler(discussionsBus *discussionsBus.Business, logger *slog.Logger) *Handler {
	return &Handler{
		discussionsBus: discussionsBus,
		logger:         logger,
	}
}

// Create
//
// @Summary Create topic message
// @Description Creates new topic message (discussion)
// @Tags topicMessages
// @Param Authorization header string true "JWT token"
// @Param request body NewTopicMessage true "Create topic message request body"
// @Accept json
// @Produce json
// @Success 201 {object} TopicMessage
// @Failure 400 {object} errs.ErrResponce
// @Failure 401 {object} errs.ErrResponce
// @Failure 403 {object} errs.ErrResponce
// @Router /topicMessages [post]
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

	topicMessage, err := h.discussionsBus.Create(ctx, busTopicMessage)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return c.Status(http.StatusCreated).JSON(toHandlerTopicMessage(topicMessage))
}
