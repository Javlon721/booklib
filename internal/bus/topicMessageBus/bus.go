package topicmessagebus

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTopicMessageNotFound = errors.New("topic's message not found")
)

type Store interface {
	Create(context.Context, TopicMessage) (uuid.UUID, error)
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

func (bus Business) Create(ctx context.Context, ntm NewTopicMessage) (TopicMessage, error) {
	now := time.Now()

	topicMessage := TopicMessage{
		Message:     ntm.Message,
		DateCreated: now,
		DateUpdated: now,
		TopicID:     ntm.TopicID,
		UserID:      ntm.UserID,
	}

	id, err := bus.store.Create(ctx, topicMessage)

	if err != nil {
		return TopicMessage{}, fmt.Errorf("create: %w", err)
	}

	topicMessage.ID = id

	return topicMessage, nil
}
