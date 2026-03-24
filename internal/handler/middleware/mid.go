package middleware

import (
	"context"
	"fmt"

	"github.com/Javlon721/booklib/internal/types/role"
	"github.com/google/uuid"
)

type ctxKey int

const (
	userIDKey ctxKey = iota
	userRolesKey
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

func setUserRoles(ctx context.Context, roles []role.Role) context.Context {
	return context.WithValue(ctx, userRolesKey, roles)
}

func GetUserRoles(ctx context.Context) ([]role.Role, error) {
	roles, ok := ctx.Value(userRolesKey).([]role.Role)

	if !ok {
		return []role.Role{}, fmt.Errorf("user roles is not found in context")
	}

	return roles, nil
}
