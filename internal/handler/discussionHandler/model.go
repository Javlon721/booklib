package discussionHandler

import (
	"fmt"
	"time"

	discussionsBus "github.com/Javlon721/booklib/internal/bus/discussionBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/google/uuid"
)

type NewTopicMessage struct {
	Message string `json:"message"`
	TopicID string `json:"topic_id"`
}

type TopicMessage struct {
	ID          string `json:"id"`
	Message     string `json:"message"`
	DateCreated string `json:"date_created"`
	DateUpdated string `json:"date_updated"`
	UserID      string `json:"user_id"`
	TopicID     string `json:"topic_id"`
}

func toBusTopicMessage(topicMessage NewTopicMessage, userID uuid.UUID) (discussionsBus.NewTopicMessage, error) {
	var fieldErrors errs.FieldErrors

	topicID, err := uuid.Parse(topicMessage.TopicID)

	if err != nil {
		fieldErrors.Add("topicID", err)
	}

	message, err := name.Parse(topicMessage.Message)

	if err != nil {
		fieldErrors.Add("message", err)
	}

	if len(fieldErrors) > 0 {
		return discussionsBus.NewTopicMessage{}, fmt.Errorf("parse: %w", fieldErrors.ToError())
	}

	return discussionsBus.NewTopicMessage{
		TopicID: topicID,
		UserID:  userID,
		Message: message,
	}, nil
}

func toHandlerTopicMessage(topicMessage discussionsBus.TopicMessage) TopicMessage {
	return TopicMessage{
		ID:          topicMessage.ID.String(),
		Message:     topicMessage.Message.String(),
		DateCreated: topicMessage.DateCreated.Format(time.RFC3339),
		DateUpdated: topicMessage.DateUpdated.Format(time.RFC3339),
		UserID:      topicMessage.UserID.String(),
		TopicID:     topicMessage.TopicID.String(),
	}
}
