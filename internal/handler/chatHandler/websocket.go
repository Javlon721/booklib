package chatHandler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/google/uuid"
)

type WsHandler struct {
	shards []*Shard
	c      int
}

func NewWsHandler(c int) *WsHandler {
	shards := make([]*Shard, 0, c)

	for range c {
		shards = append(shards, NewShard())
	}

	return &WsHandler{
		shards: shards,
		c:      c,
	}
}

func (ws *WsHandler) Register(chatID uuid.UUID, conn *websocket.Conn) (chan<- string, error) {
	shard, err := ws.GetShardBy(chatID)

	if err != nil {
		return nil, err
	}

	return shard.Register(chatID, conn), nil
}

func (ws *WsHandler) UnRegister(chatID uuid.UUID, conn *websocket.Conn) error {
	shard, err := ws.GetShardBy(chatID)

	if err != nil {
		return err
	}

	shard.UnRegister(chatID, conn)

	return nil
}

func (ws *WsHandler) GetShardBy(chatID uuid.UUID) (*Shard, error) {
	h, err := hasher(chatID.String())

	if err != nil {
		return nil, fmt.Errorf("getting conn hash: %w", err)
	}

	pos := h % ws.c

	return ws.shards[pos], nil
}

type Shard struct {
	chats map[uuid.UUID]*Group
	sync.RWMutex
}

func NewShard() *Shard {
	return &Shard{
		chats: map[uuid.UUID]*Group{},
	}
}

func (h *Shard) GetChatByID(chatID uuid.UUID) (*Group, bool) {
	h.RLock()
	defer h.RUnlock()

	chat, ok := h.chats[chatID]

	return chat, ok
}

func (h *Shard) Register(chatID uuid.UUID, conn *websocket.Conn) chan<- string {
	chat, ok := h.GetChatByID(chatID)

	h.Lock()

	if !ok {
		chat = NewGroup()

		go chat.Listen()

		h.chats[chatID] = chat
	}

	h.Unlock()

	chat.Register(conn)

	return chat.broadcastCh
}

func (h *Shard) UnRegister(chatID uuid.UUID, conn *websocket.Conn) {
	chat, ok := h.GetChatByID(chatID)

	if !ok {
		return
	}

	h.Lock()
	defer h.Unlock()

	chat.UnRegister(conn)

	if chat.Length() == 0 {
		chat.CloseBroadcast()

		delete(h.chats, chatID)
	}
}

func (h *Shard) ChatStats(chatID uuid.UUID) string {
	chat, ok := h.GetChatByID(chatID)

	if !ok {
		return "no stats"
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "ChatID: %s\n", chatID)

	chat.Stats(&builder)

	return builder.String()
}
