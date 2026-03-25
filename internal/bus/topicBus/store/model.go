package store

import (
	"fmt"
	"time"

	topicbus "github.com/Javlon721/booklib/internal/bus/topicBus"
	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/google/uuid"
)

type topicDB struct {
	ID          uuid.UUID
	Title       string
	DateCreated time.Time
	DateUpdated time.Time
	CreatedBy   uuid.UUID
}

func toDBTopic(t topicbus.Topic) topicDB {
	return topicDB{
		ID:          t.ID,
		Title:       t.Title.String(),
		DateCreated: t.DateCreated.UTC(),
		DateUpdated: t.DateUpdated.UTC(),
		CreatedBy:   t.CreatedBy,
	}
}

func toBusTopic(t topicDB) (topicbus.Topic, error) {
	title, err := name.Parse(t.Title)

	if err != nil {
		return topicbus.Topic{}, fmt.Errorf("parse title: %w", err)
	}

	return topicbus.Topic{
		ID:          t.ID,
		Title:       title,
		DateCreated: t.DateCreated.In(time.Local),
		DateUpdated: t.DateUpdated.In(time.Local),
		CreatedBy:   t.CreatedBy,
	}, nil
}
