package chatWsHandler

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/Javlon721/booklib/internal/handler/middleware"
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

// -------------------------------------------------------------------------
type Client struct {
	userID    uuid.UUID
	isClosing bool
	sync.Mutex
}

type Group struct {
	connected   map[*websocket.Conn]*Client
	broadcastCh chan string
	isActive    bool
	isClosing   bool
}

func (g *Group) Listen() error {
	if g.isActive {
		return fmt.Errorf("group is already listening")
	}

	g.isActive = true

	go func() {
		for message := range g.broadcastCh {
			data := []byte(message)

			for v := range g.connected {
				go func(conn *websocket.Conn, client *Client) {
					client.Lock()
					defer client.Unlock()

					if client.isClosing {
						return
					}

					if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
						client.isClosing = true

						conn.WriteMessage(websocket.CloseMessage, []byte{})

						conn.Close()

						g.UnRegister(conn)
					}

				}(v, g.connected[v])
			}
		}
	}()

	return nil
}

func (g *Group) Register(conn *websocket.Conn) {
	ctx, ok := conn.Locals("context").(context.Context)

	if !ok {
		panic(fmt.Errorf("websocket getting context: %w", fmt.Errorf("there is no context")))
	}

	userID, err := middleware.GetUserID(ctx)

	if err != nil {
		panic(fmt.Errorf("websocket getting userID: %w", err))
	}

	g.connected[conn] = &Client{
		userID: userID,
	}
}

func (g *Group) UnRegister(conn *websocket.Conn) {
	delete(g.connected, conn)
}

func (g *Group) Length() int {
	return len(g.connected)
}

func (g *Group) CloseBroadcast() {
	if g.isClosing {
		return
	}

	var wg sync.WaitGroup
	g.isClosing = true

	for conn, client := range g.connected {
		wg.Go(func() {
			client.Lock()
			defer client.Unlock()

			client.isClosing = true

			g.UnRegister(conn)
		})
	}

	close(g.broadcastCh)
}

func (g *Group) Broadcast(data string) {
	g.broadcastCh <- data
}

func (g *Group) Stats(report io.Writer) {
	fmt.Fprintf(report, "\tConnected %d users\n", g.Length())
	fmt.Fprintf(report, "------------------------------\n")

	for _, client := range g.connected {
		fmt.Fprintf(report, "\tuserID: %s\n", client.userID)

		var status string

		if client.isClosing {
			status = "connection closed"
		} else {
			status = "connection alive"
		}

		fmt.Fprintf(report, "\tstatus: %s\n", status)
		fmt.Fprintf(report, "------------------------------\n")
	}
}

func NewGroup() *Group {
	return &Group{
		connected:   map[*websocket.Conn]*Client{},
		broadcastCh: make(chan string),
		isActive:    false,
	}
}
