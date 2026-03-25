package topicbus

import (
	"time"

	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/google/uuid"
)

type NewTopic struct {
	Title     name.Name
	CreatedBy uuid.UUID
}

type Topic struct {
	ID          uuid.UUID
	Title       name.Name
	DateCreated time.Time
	DateUpdated time.Time
	CreatedBy   uuid.UUID
}
