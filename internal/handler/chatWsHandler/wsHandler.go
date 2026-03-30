package chatWsHandler

import (
	"fmt"

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

func (ws *WsHandler) Register(userID, chatID uuid.UUID, conn *websocket.Conn) (chan<- string, error) {
	chats, err := ws.GetChatsBy(chatID)

	if err != nil {
		return nil, err
	}

	return chats.Register(userID, chatID, conn), nil
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
