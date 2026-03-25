package topichandler

import (
	"fmt"
	"time"

	topicbus "github.com/Javlon721/booklib/internal/bus/topicBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/google/uuid"
)

type NewTopic struct {
	Title string `json:"title"`
}

func toBusTopic(t NewTopic, userID uuid.UUID) (topicbus.NewTopic, error) {
	var errors errs.FieldErrors

	title, err := name.Parse(t.Title)

	if err != nil {
		errors.Add("title", err)
	}

	if len(errors) > 0 {
		return topicbus.NewTopic{}, fmt.Errorf("validate: %w", errors.ToError())
	}

	return topicbus.NewTopic{
		Title:     title,
		CreatedBy: userID,
	}, nil
}

// -------------------------------------------------------------------------

type Topic struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	CreatedBy   string `json:"created_by"`
}

func toHandlerTopic(t topicbus.Topic) Topic {
	return Topic{
		ID:          t.ID.String(),
		Title:       t.Title.String(),
		DateCreated: t.DateCreated.Format(time.RFC3339),
		DateUpdated: t.DateUpdated.Format(time.RFC3339),
		CreatedBy:   t.CreatedBy.String(),
	}
}
