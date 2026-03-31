package chatWsHandler

import (
	"context"
	"fmt"
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

	for i := range c {
		chats = append(chats, NewChats(i))
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

func (ws *WsHandler) GetChatsByIdx(idx int) (*Chats, error) {
	idx = idx % ws.c
	return ws.chats[idx], nil
}

func (ws *WsHandler) Stats(ctx context.Context) <-chan any {
	reports := make(chan any)

	fanout := make([]<-chan any, ws.c)

	for idx, chats := range ws.chats {
		fanout[idx] = chats.Stats(ctx)
	}

	var wg sync.WaitGroup

	fanin := func(ctx context.Context, rch <-chan any) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case value, ok := <-rch:
				if !ok {
					return
				}
				reports <- value
			}
		}
	}

	wg.Add(ws.c)

	for _, rch := range fanout {
		go fanin(ctx, rch)
	}

	go func() {
		wg.Wait()
		close(reports)
	}()

	return reports
}
