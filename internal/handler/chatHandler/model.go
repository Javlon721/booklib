package chatHandler

import (
	"context"
	"fmt"
	"time"

	chatbus "github.com/Javlon721/booklib/internal/bus/chatBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/handler/middleware"
	"github.com/Javlon721/booklib/internal/types/name"
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
