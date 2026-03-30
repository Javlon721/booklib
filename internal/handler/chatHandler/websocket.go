package chatHandler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/google/uuid"
)

type WsHandler struct {
	chats []*Chats
	c     int
}

func NewWsHandler(c int) *WsHandler {
	chats := make([]*Chats, 0, c)

	for range c {
		chats = append(chats, NewChats())
	}

	return &WsHandler{
		chats: chats,
		c:     c,
	}
}

func (ws *WsHandler) Register(chatID uuid.UUID, conn *websocket.Conn) (chan<- string, error) {
	chats, err := ws.GetChatsBy(chatID)

	if err != nil {
		return nil, err
	}

	return chats.Register(chatID, conn), nil
}

func (ws *WsHandler) UnRegister(chatID uuid.UUID, conn *websocket.Conn) error {
	chats, err := ws.GetChatsBy(chatID)

	if err != nil {
		return err
	}

	chats.UnRegister(chatID, conn)

	return nil
}

func (ws *WsHandler) GetChatsBy(chatID uuid.UUID) (*Chats, error) {
	h, err := hasher(chatID.String())

	if err != nil {
		return nil, fmt.Errorf("getting conn hash: %w", err)
	}

	pos := h % ws.c

	return ws.chats[pos], nil
}

// -------------------------------------------------------------------------

type Chats struct {
	chats map[uuid.UUID]*Group
	sync.RWMutex
}

func NewChats() *Chats {
	return &Chats{
		chats: map[uuid.UUID]*Group{},
	}
}

func (h *Chats) GetChatByID(chatID uuid.UUID) (*Group, bool) {
	h.RLock()
	defer h.RUnlock()

	chat, ok := h.chats[chatID]

	return chat, ok
}

func (h *Chats) Register(chatID uuid.UUID, conn *websocket.Conn) chan<- string {
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

func (h *Chats) UnRegister(chatID uuid.UUID, conn *websocket.Conn) {
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

func (h *Chats) ChatStats(chatID uuid.UUID) string {
	chat, ok := h.GetChatByID(chatID)

	if !ok {
		return "no stats"
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "ChatID: %s\n", chatID)

	chat.Stats(&builder)

	return builder.String()
}
