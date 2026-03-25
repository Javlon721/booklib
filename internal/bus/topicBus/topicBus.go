package topicbus

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTopicNotFound      = errors.New("topic not found")
	ErrTopicAlreadyExists = errors.New("topic already exists")
)

type Store interface {
	Create(context.Context, Topic) (uuid.UUID, error)
	GetByID(context.Context, uuid.UUID) (Topic, error)
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

func (bus Business) Create(ctx context.Context, nt NewTopic) (Topic, error) {
	now := time.Now()

	topic := Topic{
		Title:       nt.Title,
		DateCreated: now,
		DateUpdated: now,
		CreatedBy:   nt.CreatedBy,
	}

	id, err := bus.store.Create(ctx, topic)

	if err != nil {
		return Topic{}, fmt.Errorf("create: %w", err)
	}

	topic.ID = id

	return topic, nil
}

func (bus Business) GetUserByID(ctx context.Context, topicID uuid.UUID) (Topic, error) {
	return bus.store.GetByID(ctx, topicID)
}
