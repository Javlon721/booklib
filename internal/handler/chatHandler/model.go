package chatHandler

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	chatbus "github.com/Javlon721/booklib/internal/bus/chatBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/handler/middleware"
	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/google/uuid"
)

type NewChat struct {
	Reciever string `json:"reciever"`
}

func toBusChat(chat NewChat, sender uuid.UUID) (chatbus.NewChat, error) {
	var fieldErrors errs.FieldErrors

	reciever, err := uuid.Parse(chat.Reciever)

	if err != nil {
		fieldErrors.Add("reciever", err)
	}

	if len(fieldErrors) > 0 {
		return chatbus.NewChat{}, fmt.Errorf("parse: %w", fieldErrors.ToError())
	}

	return chatbus.NewChat{
		User1: sender,
		User2: reciever,
	}, nil
}

// -------------------------------------------------------------------------

type Chat struct {
	ID          string `json:"chat_id"`
	User1       string `json:"user1"`
	User2       string `json:"user2"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
}

func toHandlerChat(chat chatbus.Chat) Chat {
	return Chat{
		ID:          chat.ID.String(),
		User1:       chat.User1.String(),
		User2:       chat.User2.String(),
		DateCreated: chat.DateCreated.Format(time.RFC3339),
		DateUpdated: chat.DateUpdated.Format(time.RFC3339),
	}
}

// -------------------------------------------------------------------------
type NewMessage struct {
	Content string `json:"content"`
}

func toBusMessage(ctx context.Context, m NewMessage, chatID string) (chatbus.NewMessage, error) {
	var fieldErrors errs.FieldErrors

	sender, err := middleware.GetUserID(ctx)

	if err != nil {
		fieldErrors.Add("sender", err)
	}

	id, err := uuid.Parse(chatID)

	if err != nil {
		fieldErrors.Add("chatID", err)
	}

	content, err := name.Parse(m.Content)

	if err != nil {
		fieldErrors.Add("content", err)
	}

	if len(fieldErrors) > 0 {
		return chatbus.NewMessage{}, fmt.Errorf("parse: %w", fieldErrors.ToError())
	}

	return chatbus.NewMessage{
		ChatID:  id,
		Sender:  sender,
		Content: content,
	}, nil
}

// -------------------------------------------------------------------------

type Message struct {
	ID          string `json:"id"`
	ChatID      string `json:"chat_id"`
	Sender      string `json:"sender"`
	Content     string `json:"content"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	Status      string `json:"status"`
}

func toHandlerMessage(m chatbus.Message) Message {
	return Message{
		ID:          m.ID.String(),
		ChatID:      m.ChatID.String(),
		Sender:      m.Sender.String(),
		Content:     m.Content.String(),
		Status:      m.Status.String(),
		DateCreated: m.DateCreated.Format(time.RFC3339),
		DateUpdated: m.DateUpdated.Format(time.RFC3339),
	}
}

func toHandlerMessages(messages []chatbus.Message) []Message {
	result := make([]Message, len(messages))

	for i, v := range messages {
		result[i] = toHandlerMessage(v)
	}

	return result
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
