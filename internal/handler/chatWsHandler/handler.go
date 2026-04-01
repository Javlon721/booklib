package chatWsHandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	chatbus "github.com/Javlon721/booklib/internal/bus/chatBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/handler/middleware"
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

// Websoket
//
// @Summary Websocket connection
// @Description Websocket connection for chating
// @Tags ws/chats
// @Param Authorization header string true "JWT token"
// @Param chatID path string true "Chat uuid"
// @Produce plain
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} errs.ErrResponce
// @Failure 426 {object} errs.ErrResponce
// @Router /ws/chats/{shatID} [get]
func (h Handler) Websoket(c *websocket.Conn) {
	defer c.Close()

	ctx, ok := c.Locals("context").(context.Context)

	if !ok {
		h.logger.Error("websocket getting context", "err", errs.New(errs.Internal, fmt.Errorf("there is no context")))
		return
	}

	userID, err := middleware.GetUserID(ctx)

	if err != nil {
		c.WriteMessage(websocket.CloseMessage, []byte(fmt.Errorf("websocket getting userID: %w", err).Error()))
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

	broadcastCh, err := h.wsHandler.Register(userID, chatID, c)

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

// -------------------------------------------------------------------------
// Stats

// WebsoketChatStats
//
// @Summary Get ws chat stats
// @Description Gets Websocket Chat stats
// @Tags ws/chats
// @Param chatID path string true "Chat uuid"
// @Produce json
// @Success 200 {object} []ChatStats
// @Failure 400 {object} errs.ErrResponce
// @Failure 404 {object} errs.ErrResponce
// @Router /ws/chats/stats/chat/{chatID} [get]
func (h Handler) WebsoketChatStats(c fiber.Ctx) error {
	chatID, err := uuid.Parse(c.Params("chatID"))

	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	chats, err := h.wsHandler.GetChatsBy(chatID)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	stats, ok := chats.ChatStats(chatID)

	if !ok {
		return errs.New(errs.NotFound, fmt.Errorf("no stats"))
	}

	return c.JSON(stats)
}

// WebsoketChatsStats
//
// @Summary Get ws chats stats
// @Description Gets Websocket Chats stats
// @Tags ws/chats
// @Param idx path int true "Shard id"
// @Produce application/x-ndjson
// @Success 200 {object} []ShardStats
// @Failure 400 {object} errs.ErrResponce
// @Router /ws/chats/stats/{idx} [get]
func (h Handler) WebsoketChatsStats(c fiber.Ctx) error {
	idx, err := strconv.Atoi(c.Params("idx"))

	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	chats, err := h.wsHandler.GetChatsByIdx(idx)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	ctx := c.Context()

	reports := chats.Stats(ctx)

	reader := NewChanReader(ctx)

	go reader.Listen(reports)

	c.Response().Header.Add("Content-Type", "application/x-ndjson")

	return c.SendStream(reader)
}

// WebsoketShardsStats
//
// @Summary Get ws shards stats
// @Description Gets Websocket shards stats
// @Tags ws/chats
// @Produce application/x-ndjson
// @Success 200 {object} []ShardStats
// @Failure 400 {object} errs.ErrResponce
// @Router /ws/chats/stats [get]
func (h Handler) WebsoketShardsStats(c fiber.Ctx) error {
	ctx := c.Context()

	reports := h.wsHandler.Stats(ctx)
	reader := NewChanReader(ctx)

	go reader.Listen(reports)

	c.Response().Header.Add("Content-Type", "application/x-ndjson")
	return c.SendStream(reader)
}
