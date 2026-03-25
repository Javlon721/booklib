package store

import (
	"fmt"
	"time"

	topicmessagebus "github.com/Javlon721/booklib/internal/bus/topicMessageBus"
	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/google/uuid"
)

type topicMessageDB struct {
	ID          uuid.UUID `db:"id"`
	Message     string    `db:"message"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
	UserID      uuid.UUID `db:"user_id"`
	TopicID     uuid.UUID `db:"topic_id"`
}

func toDBTopicMessage(m topicmessagebus.TopicMessage) topicMessageDB {
	return topicMessageDB{
		ID:          m.ID,
		Message:     m.Message.String(),
		DateCreated: m.DateCreated.UTC(),
		DateUpdated: m.DateUpdated.UTC(),
		UserID:      m.UserID,
		TopicID:     m.TopicID,
	}
}

func toBusTopicMessage(m topicMessageDB) (topicmessagebus.TopicMessage, error) {
	message, err := name.Parse(m.Message)

	if err != nil {
		return topicmessagebus.TopicMessage{}, fmt.Errorf("parse message: %w", err)
	}

	return topicmessagebus.TopicMessage{
		ID:          m.ID,
		Message:     message,
		DateCreated: m.DateCreated.In(time.Local),
		DateUpdated: m.DateUpdated.In(time.Local),
		UserID:      m.UserID,
		TopicID:     m.TopicID,
	}, nil
}

func toBusTopicMessages(messages []topicMessageDB) ([]topicmessagebus.TopicMessage, error) {
	result := make([]topicmessagebus.TopicMessage, len(messages))

	for i, message := range messages {
		busMessage, err := toBusTopicMessage(message)

		if err != nil {
			return []topicmessagebus.TopicMessage{}, err
		}

		result[i] = busMessage
	}

	return result, nil
}
