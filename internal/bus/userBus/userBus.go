package userbus

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	Create(context.Context, User) (uuid.UUID, error)
	GetUserByID(context.Context, uuid.UUID) (User, error)
}

type Business struct {
	logger *slog.Logger
	store  Store
}

func NewBusiness(logger *slog.Logger, store Store) *Business {
	return &Business{
		store:  store,
		logger: logger,
	}
}

func (bus Business) Create(ctx context.Context, nu CreateUser) (User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(nu.Password.String()), bcrypt.DefaultCost)

	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	user := User{
		FirstName:    nu.FirstName,
		LastName:     nu.LastName,
		PasswordHash: passwordHash,
		Email:        nu.Email,
		DateCreated:  now,
		DateUpdated:  now,
	}

	id, err := bus.store.Create(ctx, user)

	if err != nil {
		return User{}, fmt.Errorf("create: %w", err)
	}

	user.ID = id

	return user, nil
}

func (bus Business) GetUserByID(ctx context.Context, userID uuid.UUID) (User, error) {
	return bus.store.GetUserByID(ctx, userID)
}
