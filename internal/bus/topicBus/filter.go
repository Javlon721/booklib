package topicbus

import (
	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/google/uuid"
)

type QueryFilter struct {
	CreatedBy *uuid.UUID
	Title     *name.Name
}
