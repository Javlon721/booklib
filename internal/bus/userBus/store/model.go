package store

import (
	"fmt"
	"net/mail"
	"time"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/Javlon721/booklib/internal/types/name"
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
		Email:        u.Email.Address,
		PasswordHash: u.PasswordHash,
		DateCreated:  u.DateCreated.UTC(),
		DateUpdated:  u.DateUpdated.UTC(),
	}
}

func toBusUser(u userDB) (userbus.User, error) {
	email := mail.Address{
		Address: u.Email,
	}

	firstName, err := name.Parse(u.FirstName)

	if err != nil {
		return userbus.User{}, fmt.Errorf("parse firstName: %w", err)
	}

	lastName, err := name.Parse(u.LastName)

	if err != nil {
		return userbus.User{}, fmt.Errorf("parse lastName: %w", err)
	}

	return userbus.User{
		ID:           u.ID,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: u.PasswordHash,
		DateCreated:  u.DateCreated.In(time.Local),
		DateUpdated:  u.DateUpdated.In(time.Local),
	}, nil
}
