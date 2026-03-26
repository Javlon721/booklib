package store

import (
	"fmt"
	"time"

	discussionBus "github.com/Javlon721/booklib/internal/bus/discussionBus"
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

func toDBTopicMessage(m discussionBus.TopicMessage) topicMessageDB {
	return topicMessageDB{
		ID:          m.ID,
		Message:     m.Message.String(),
		DateCreated: m.DateCreated.UTC(),
		DateUpdated: m.DateUpdated.UTC(),
		UserID:      m.UserID,
		TopicID:     m.TopicID,
	}
}

func toBusTopicMessage(m topicMessageDB) (discussionBus.TopicMessage, error) {
	message, err := name.Parse(m.Message)

	if err != nil {
		return discussionBus.TopicMessage{}, fmt.Errorf("parse message: %w", err)
	}

	return discussionBus.TopicMessage{
		ID:          m.ID,
		Message:     message,
		DateCreated: m.DateCreated.In(time.Local),
		DateUpdated: m.DateUpdated.In(time.Local),
		UserID:      m.UserID,
		TopicID:     m.TopicID,
	}, nil
}

func toBusTopicMessages(messages []topicMessageDB) ([]discussionBus.TopicMessage, error) {
	result := make([]discussionBus.TopicMessage, len(messages))

	for i, message := range messages {
		busMessage, err := toBusTopicMessage(message)

		if err != nil {
			return []discussionBus.TopicMessage{}, err
		}

		result[i] = busMessage
	}

	return result, nil
}
