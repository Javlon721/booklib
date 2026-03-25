package store

import (
	"fmt"
	"time"

	topicbus "github.com/Javlon721/booklib/internal/bus/topicBus"
	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/google/uuid"
)

type topicDB struct {
	ID          uuid.UUID `db:"topic_id"`
	Title       string    `db:"title"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
	CreatedBy   uuid.UUID `db:"created_by"`
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

func toBusTopics(t []topicDB) ([]topicbus.Topic, error) {
	result := make([]topicbus.Topic, len(t))

	for i, v := range t {
		busTopic, err := toBusTopic(v)

		if err != nil {
			return nil, err
		}

		result[i] = busTopic
	}

	return result, nil
}
