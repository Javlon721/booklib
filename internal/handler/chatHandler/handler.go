package chatHandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	chatbus "github.com/Javlon721/booklib/internal/bus/chatBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/handler/middleware"
	"github.com/Javlon721/booklib/internal/types/messageStatus"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler struct {
	chatBus   *chatbus.Business
	logger    *slog.Logger
	wsHandler *WsHandler
}

func NewHandler(chatBus *chatbus.Business, logger *slog.Logger) *Handler {
	return &Handler{
		chatBus:   chatBus,
		logger:    logger,
		wsHandler: NewWsHandler(4),
	}
}

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

func (h Handler) Websoket(c *websocket.Conn) {
	defer c.Close()

	ctx, ok := c.Locals("context").(context.Context)

	if !ok {
		h.logger.Error("websocket getting context", "err", errs.New(errs.Internal, fmt.Errorf("there is no context")))
		return
	}

	chatID, err := uuid.Parse(c.Params("chatID"))

	if err != nil {
		c.WriteMessage(websocket.CloseMessage, []byte(errs.InvalidArgument.String()))
		return
	}

	_, err = h.chatBus.GetChatByID(ctx, chatID)

	if err != nil {
		var message error

		if errors.Is(err, chatbus.ErrChatNotFound) {
			message = err
		} else {
			message = fmt.Errorf("some error occured")
		}

		h.logger.Error("websocket retrieve chat", "err", err)

		c.WriteMessage(websocket.CloseMessage, []byte(message.Error()))

		return
	}

	broadcastCh, err := h.wsHandler.Register(chatID, c)

	if err != nil {
		h.logger.Error("websocket register chat", "err", err)
		c.WriteMessage(websocket.CloseMessage, []byte(err.Error()))
		return
	}

	defer h.wsHandler.UnRegister(chatID, c)

	for {
		mt, msg, err := c.ReadMessage()

		if err != nil {
			h.logger.Error("websocket read", "err", err)

			c.WriteMessage(websocket.CloseMessage, []byte("unexpected error while reading"))

			return
		}

		if mt != websocket.TextMessage {
			resp := "wrong message type"

			h.logger.Error("websocket message type", "err", resp)

			c.WriteMessage(websocket.TextMessage, []byte(resp))

			continue
		}

		broadcastCh <- string(msg)
	}
}

func (h Handler) WebsoketStats(c fiber.Ctx) error {
	chatID, err := uuid.Parse(c.Params("chatID"))

	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	chats, err := h.wsHandler.GetChatsBy(chatID)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return c.SendString(chats.ChatStats(chatID))
}
