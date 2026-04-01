package chatHandler

import (
	"errors"
	"log/slog"
	"net/http"

	chatbus "github.com/Javlon721/booklib/internal/bus/chatBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/handler/middleware"
	"github.com/Javlon721/booklib/internal/types/messageStatus"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler struct {
	chatBus *chatbus.Business
	logger  *slog.Logger
}

func NewHandler(chatBus *chatbus.Business, logger *slog.Logger) *Handler {
	return &Handler{
		chatBus: chatBus,
		logger:  logger,
	}
}

// Create
//
// @Summary Create chat
// @Description Creates new chat between users
// @Tags chats
// @Param Authorization header string true "JWT token"
// @Param request body NewChat true "Create chat request body (sender already included via JWT token)"
// @Accept json
// @Produce json
// @Success 201 {object} Chat
// @Failure 400 {object} errs.ErrResponce
// @Failure 401 {object} errs.ErrResponce
// @Failure 403 {object} errs.ErrResponce
// @Failure 409 {object} errs.ErrResponce
// @Router /chats/chat [post]
func (h Handler) CreateChat(c fiber.Ctx) error {
	ctx := c.Context()

	sender, err := middleware.GetUserID(ctx)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	var newChat NewChat

	if err := c.Bind().Body(&newChat); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	busChat, err := toBusChat(newChat, sender)

	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	chat, err := h.chatBus.CreateChat(ctx, busChat)

	if err != nil {
		if errors.Is(err, chatbus.ErrChatUsersViolation) {
			return errs.New(errs.InvalidArgument, err)
		}
		if errors.Is(err, chatbus.ErrChatAlreadyExists) {
			return errs.New(errs.AlreadyExists, err)
		}
		return errs.New(errs.Internal, err)
	}

	return c.Status(http.StatusCreated).JSON(toHandlerChat(chat))
}

// GetChatBy
//
// @Summary Get chat
// @Description Get chat by reciever (sender already included via JWT token)
// @Tags chats
// @Param Authorization header string true "JWT token"
// @Param userID path string true "Reciever uuid"
// @Produce json
// @Success 200 {object} Chat
// @Failure 400 {object} errs.ErrResponce
// @Failure 401 {object} errs.ErrResponce
// @Failure 403 {object} errs.ErrResponce
// @Router /chats/chat/{userID} [get]
func (h Handler) GetChatBy(c fiber.Ctx) error {
	ctx := c.Context()

	sender, err := middleware.GetUserID(ctx)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	reciever, err := uuid.Parse(c.Params("userID"))

	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	chat, err := h.chatBus.GetChatBy(ctx, sender, reciever)

	if err != nil {
		if errors.Is(err, chatbus.ErrChatNotFound) {
			return errs.New(errs.NotFound, err)
		}

		if errors.Is(err, chatbus.ErrChatUsersViolation) {
			return errs.New(errs.Aborted, err)
		}

		return errs.New(errs.Internal, err)
	}

	return c.JSON(toHandlerChat(chat))
}

// CreateMessage
//
// @Summary Create new message
// @Description Creates new message between users
// @Tags chats
// @Param Authorization header string true "JWT token"
// @Param request body NewMessage true "Create new message request body"
// @Param chatID path string true "Chat id for creating new message"
// @Accept json
// @Produce json
// @Success 201 {object} Message
// @Failure 400 {object} errs.ErrResponce
// @Failure 401 {object} errs.ErrResponce
// @Failure 403 {object} errs.ErrResponce
// @Failure 404 {object} errs.ErrResponce
// @Router /chats/chat/{chatID} [post]
func (h Handler) CreateMessage(c fiber.Ctx) error {
	var newMessage NewMessage

	if err := c.Bind().Body(&newMessage); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ctx := c.Context()

	// todo if chat does not exist, now yeids err: foregn key violation
	chatID := c.Params("chatID")

	busNewMessage, err := toBusMessage(ctx, newMessage, chatID)

	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	message, err := h.chatBus.CreateMessage(ctx, busNewMessage)

	if err != nil {
		if errors.Is(err, chatbus.ErrChatNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return errs.New(errs.Internal, err)
	}

	return c.Status(http.StatusCreated).JSON(toHandlerMessage(message))
}

// GetPendingMessages
//
// @Summary Get pending messages
// @Description Gets only pending messages
// @Tags chats
// @Param Authorization header string true "JWT token"
// @Param chatID path string true "Chat uuid"
// @Produce json
// @Success 200 {object} []Message
// @Failure 400 {object} errs.ErrResponce
// @Failure 401 {object} errs.ErrResponce
// @Failure 403 {object} errs.ErrResponce
// @Router /chats/messages/{chatID} [get]
func (h Handler) GetPendingMessages(c fiber.Ctx) error {
	ctx := c.Context()

	reciever, err := middleware.GetUserID(ctx)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	chatID, err := uuid.Parse(c.Params("chatID"))

	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	messages, err := h.chatBus.GetMessagesByStatus(ctx, reciever, chatID, messageStatus.Pending)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return c.JSON(toHandlerMessages(messages))
}
