package middleware

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ctxKey int

const (
	userIDKey ctxKey = iota
)

func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(userIDKey).(uuid.UUID)

	if !ok {
		return uuid.UUID{}, fmt.Errorf("userID not found in context")
	}

	return userID, nil
}
