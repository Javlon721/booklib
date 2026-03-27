package chatbus

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Javlon721/booklib/internal/types/messageStatus"
	"github.com/google/uuid"
)

var (
	ErrChatNotFound       = errors.New("chat not found")
	ErrChatAlreadyExists  = errors.New("chat already exists")
	ErrChatUsersViolation = errors.New("users id should not be same")
)

type Store interface {
	CreateChat(context.Context, Chat) (uuid.UUID, error)
	GetChatBy(context.Context, uuid.UUID, uuid.UUID) (Chat, error)
	CreateMessage(context.Context, Message) (uuid.UUID, error)
	GetMessagesByStatus(context.Context, uuid.UUID, uuid.UUID, string) ([]Message, error)
}

type Business struct {
	logger *slog.Logger
	store  Store
}

func NewBusiness(logger *slog.Logger, store Store) *Business {
	return &Business{
		store:  store,
		logger: logger,
	}
}

func (bus Business) CreateChat(ctx context.Context, nch NewChat) (Chat, error) {
	if nch.User1 == nch.User2 {
		return Chat{}, ErrChatUsersViolation
	}

	now := time.Now()

	chat := Chat{
		DateCreated: now,
		DateUpdated: now,
		User1:       nch.User1,
		User2:       nch.User2,
	}

	chatID, err := bus.store.CreateChat(ctx, chat)

	if err != nil {
		return Chat{}, err
	}

	chat.ID = chatID

	return chat, nil
}

func (bus Business) GetChatBy(ctx context.Context, user1, user2 uuid.UUID) (Chat, error) {
	if user1 == user2 {
		return Chat{}, ErrChatUsersViolation
	}

	return bus.store.GetChatBy(ctx, user1, user2)
}

func (bus Business) CreateMessage(ctx context.Context, nm NewMessage) (Message, error) {
	now := time.Now()

	status := messageStatus.Pending

	message := Message{
		ChatID:      nm.ChatID,
		Sender:      nm.Sender,
		DateCreated: now,
		DateUpdated: now,
		Content:     nm.Content,
		Status:      status,
	}

	messageID, err := bus.store.CreateMessage(ctx, message)

	if err != nil {
		return Message{}, err
	}

	message.ID = messageID

	return message, nil
}

func (bus Business) GetMessagesByStatus(
	ctx context.Context, reciever, chatID uuid.UUID, status messageStatus.Status) ([]Message, error) {
	return bus.store.GetMessagesByStatus(ctx, chatID, reciever, status.String())
}
