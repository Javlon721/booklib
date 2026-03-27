package chatHandler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/google/uuid"
)

type WebsoketHandler struct {
	chats map[uuid.UUID]*Group
	sync.RWMutex
}

func NewWebsoketHandler() *WebsoketHandler {
	return &WebsoketHandler{
		chats: map[uuid.UUID]*Group{},
	}
}

func (h *WebsoketHandler) GetChatByID(chatID uuid.UUID) (*Group, bool) {
	h.RLock()
	defer h.RUnlock()

	chat, ok := h.chats[chatID]

	return chat, ok
}

func (h *WebsoketHandler) Register(chatID uuid.UUID, conn *websocket.Conn) chan<- string {
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

func (h *WebsoketHandler) UnRegister(chatID uuid.UUID, conn *websocket.Conn) {
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

func (h *WebsoketHandler) ChatStats(chatID uuid.UUID) string {
	chat, ok := h.GetChatByID(chatID)

	if !ok {
		return "no stats"
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "ChatID: %s\n", chatID)

	chat.Stats(&builder)

	return builder.String()
}
