package chatbus

import (
	"time"

	"github.com/Javlon721/booklib/internal/types/messageStatus"
	"github.com/Javlon721/booklib/internal/types/name"
	"github.com/google/uuid"
)

type NewChat struct {
	User1 uuid.UUID
	User2 uuid.UUID
}

type Chat struct {
	ID          uuid.UUID
	User1       uuid.UUID
	User2       uuid.UUID
	DateCreated time.Time
	DateUpdated time.Time
}

type NewMessage struct {
	ChatID  uuid.UUID
	Sender  uuid.UUID
	Content name.Name
}

type Message struct {
	ID          uuid.UUID
	ChatID      uuid.UUID
	Sender      uuid.UUID
	Content     name.Name
	DateCreated time.Time
	DateUpdated time.Time
	Status      messageStatus.Status
}
