package chatWsHandler

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/google/uuid"
)

type Client struct {
	userID    uuid.UUID
	isClosing bool
	sync.Mutex
}

// -------------------------------------------------------------------------

type Chat struct {
	connected   map[*websocket.Conn]*Client
	broadcastCh chan string
	isActive    bool
	isClosing   bool
}

func (c *Chat) Listen() error {
	if c.isActive {
		return fmt.Errorf("chat is already listening")
	}

	c.isActive = true

	go func() {
		for message := range c.broadcastCh {
			data := []byte(message)

			for v := range c.connected {
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

						c.UnRegister(conn)
					}

				}(v, c.connected[v])
			}
		}
	}()

	return nil
}

func (c *Chat) Register(userID uuid.UUID, conn *websocket.Conn) {
	c.connected[conn] = &Client{
		userID: userID,
	}
}

func (c *Chat) UnRegister(conn *websocket.Conn) {
	delete(c.connected, conn)
}

func (c *Chat) Length() int {
	return len(c.connected)
}

func (c *Chat) CloseBroadcast() {
	if c.isClosing {
		return
	}

	var wg sync.WaitGroup
	c.isClosing = true

	for conn, client := range c.connected {
		wg.Go(func() {
			client.Lock()
			defer client.Unlock()

			client.isClosing = true

			c.UnRegister(conn)
		})
	}

	close(c.broadcastCh)
}

func (c *Chat) Broadcast(data string) {
	c.broadcastCh <- data
}

func (c *Chat) Stats(report io.Writer) {
	fmt.Fprintf(report, "\tConnected %d users\n", c.Length())
	fmt.Fprintf(report, "------------------------------\n")

	for _, client := range c.connected {
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

func NewChat() *Chat {
	return &Chat{
		connected:   map[*websocket.Conn]*Client{},
		broadcastCh: make(chan string),
		isActive:    false,
	}
}

// -------------------------------------------------------------------------

type Chats struct {
	chats map[uuid.UUID]*Chat
	sync.RWMutex
}

func NewChats() *Chats {
	return &Chats{
		chats: map[uuid.UUID]*Chat{},
	}
}

func (h *Chats) GetChatByID(chatID uuid.UUID) (*Chat, bool) {
	h.RLock()
	defer h.RUnlock()

	chat, ok := h.chats[chatID]

	return chat, ok
}

func (h *Chats) Register(userID, chatID uuid.UUID, conn *websocket.Conn) chan<- string {
	chat, ok := h.GetChatByID(chatID)

	h.Lock()

	if !ok {
		chat = NewChat()

		go chat.Listen()

		h.chats[chatID] = chat
	}

	h.Unlock()

	chat.Register(userID, conn)

	return chat.broadcastCh
}

// todo something is off in unregistering. I'm not satisfied with mutexes somehow
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
