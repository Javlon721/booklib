package store

import (
	"fmt"
	"time"

	chatbus "github.com/Javlon721/booklib/internal/bus/chatBus"
	"github.com/Javlon721/booklib/internal/types/messageStatus"
	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/google/uuid"
)

type chatDB struct {
	ID          uuid.UUID `db:"chat_id"`
	User1       uuid.UUID `db:"user1"`
	User2       uuid.UUID `db:"user2"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBChat(chat chatbus.Chat) chatDB {
	return chatDB{
		ID:          chat.ID,
		User1:       chat.User1,
		User2:       chat.User2,
		DateCreated: chat.DateCreated.UTC(),
		DateUpdated: chat.DateUpdated.UTC(),
	}
}

func toBusChat(chat chatDB) (chatbus.Chat, error) {
	return chatbus.Chat{
		ID:          chat.ID,
		User1:       chat.User1,
		User2:       chat.User2,
		DateCreated: chat.DateCreated.In(time.Local),
		DateUpdated: chat.DateUpdated.In(time.Local),
	}, nil
}

// -------------------------------------------------------------------------

type messageDB struct {
	ID          uuid.UUID `db:"id"`
	ChatID      uuid.UUID `db:"chat_id"`
	Sender      uuid.UUID `db:"sender"`
	Content     string    `db:"content"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
	Status      string    `db:"status"`
}

func toDBmessage(m chatbus.Message) messageDB {
	return messageDB{
		ID:          m.ID,
		ChatID:      m.ChatID,
		Sender:      m.Sender,
		Content:     m.Content.String(),
		DateCreated: m.DateCreated.UTC(),
		DateUpdated: m.DateUpdated.UTC(),
		Status:      m.Status.String(),
	}
}

func toBusMessage(m messageDB) (chatbus.Message, error) {
	content, err := name.Parse(m.Content)

	if err != nil {
		return chatbus.Message{}, fmt.Errorf("parse content: %w", err)
	}

	status, err := messageStatus.Parse(m.Status)

	if err != nil {
		return chatbus.Message{}, fmt.Errorf("parse status: %w", err)
	}

	return chatbus.Message{
		ID:          m.ID,
		ChatID:      m.ChatID,
		Sender:      m.Sender,
		Content:     content,
		DateCreated: m.DateCreated.In(time.Local),
		DateUpdated: m.DateUpdated.In(time.Local),
		Status:      status,
	}, nil
}

func toBusMessages(messages []messageDB) ([]chatbus.Message, error) {
	result := make([]chatbus.Message, len(messages))

	for i, v := range messages {
		message, err := toBusMessage(v)

		if err != nil {
			return nil, err
		}

		result[i] = message
	}

	return result, nil
}
