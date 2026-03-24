package authbus

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/Javlon721/booklib/internal/types/role"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserIDMissing      = errors.New("userID missing from token hander")
	ErrUserRolesMissing   = errors.New("user roles missing from token hander")
	ErrUserIDMalformed    = errors.New("userID in token header is malformed")
	ErrUserRolesMalformed = errors.New("user roles in token header are malformed")
	ErrInvalidToken       = errors.New("invalid token")
)

type TokenConfig struct {
	Secret         []byte
	TokenExpiresAt time.Duration
	Method         jwt.SigningMethod
}

type UserBus interface {
	GetUserByID(context.Context, uuid.UUID) (userbus.User, error)
	GetUserByEmail(context.Context, string) (userbus.User, error)
}

type Bussiness struct {
	logger   *slog.Logger
	userBus  UserBus
	tokenCfg TokenConfig
	parser   *jwt.Parser
}

func NewBussiness(logger *slog.Logger, userBus UserBus, tokenCfg TokenConfig) *Bussiness {
	return &Bussiness{
		logger:   logger,
		userBus:  userBus,
		tokenCfg: tokenCfg,
		parser:   jwt.NewParser(jwt.WithValidMethods([]string{tokenCfg.Method.Alg()})),
	}
}

func (bus Bussiness) Login(ctx context.Context, payload AuthParams) (string, error) {
	user, err := bus.userBus.GetUserByEmail(ctx, payload.Email.Address)

	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password.String())); err != nil {
		return "", ErrInvalidCredentials
	}

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(bus.tokenCfg.TokenExpiresAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: role.ParseToString(user.Roles),
	}

	token, err := GenerateToken(claims, bus.tokenCfg.Secret, bus.tokenCfg.Method)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (bus Bussiness) Authenticate(ctx context.Context, bearerString string) (AuthenticateResp, error) {
	if !strings.HasPrefix(bearerString, "Bearer ") {
		return AuthenticateResp{}, fmt.Errorf("expected authorization header format: Bearer <token>")
	}

	jwtUnverified := bearerString[7:]

	var claims Claims

	_, err := bus.parser.ParseWithClaims(jwtUnverified, &claims, func(t *jwt.Token) (any, error) {
		return bus.tokenCfg.Secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrSignatureInvalid) {
			return AuthenticateResp{}, ErrInvalidToken
		}
		return AuthenticateResp{}, err
	}

	if claims.Subject == "" {
		return AuthenticateResp{}, ErrUserIDMissing
	}

	userID, err := uuid.Parse(claims.Subject)

	if err != nil {
		return AuthenticateResp{}, ErrUserIDMalformed
	}

	roles, err := role.ParseMany(claims.Roles)

	if err != nil {
		return AuthenticateResp{}, ErrUserRolesMalformed
	}

	_, err = bus.userBus.GetUserByID(ctx, userID)

	if err != nil {
		return AuthenticateResp{}, ErrInvalidCredentials
	}

	return AuthenticateResp{
		UserID: userID,
		Roles:  roles,
	}, nil
}

func (bus Bussiness) Authorize(ctx context.Context, targetRoles, roles []role.Role) bool {
	for _, role := range targetRoles {
		if slices.Index(roles, role) == -1 {
			return false
		}
	}
	return true
}
