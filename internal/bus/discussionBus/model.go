package discussionBus

import (
	"time"

	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/google/uuid"
)

type NewTopicMessage struct {
	Message name.Name
	UserID  uuid.UUID
	TopicID uuid.UUID
}

type TopicMessage struct {
	ID          uuid.UUID
	Message     name.Name
	TopicID     uuid.UUID
	DateCreated time.Time
	DateUpdated time.Time
	UserID      uuid.UUID
}
