package store

import (
	"time"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/google/uuid"
)

type userDB struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	Email        string
	PasswordHash []byte
	DateCreated  time.Time
	DateUpdated  time.Time
}

func toDBUser(u userbus.User) userDB {
	return userDB{
		ID:           u.ID,
		FirstName:    u.FirstName.String(),
		LastName:     u.LastName.String(),
		Email:        u.Email.String(),
		PasswordHash: u.PasswordHash,
		DateCreated:  u.DateCreated.UTC(),
		DateUpdated:  u.DateUpdated.UTC(),
	}
}
