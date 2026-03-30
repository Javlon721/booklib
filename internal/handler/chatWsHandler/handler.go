package chatWsHandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	chatbus "github.com/Javlon721/booklib/internal/bus/chatBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
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
