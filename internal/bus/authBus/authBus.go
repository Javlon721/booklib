package authbus

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserIDMissing      = errors.New("userID missing from token hander")
	ErrUserIDMalformed    = errors.New("userID in token header is malformed")
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
}

func NewBussiness(logger *slog.Logger, userBus UserBus, tokenCfg TokenConfig) *Bussiness {
	return &Bussiness{
		logger:   logger,
		userBus:  userBus,
		tokenCfg: tokenCfg,
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

	exp := time.Now().UTC().Add(bus.tokenCfg.TokenExpiresAt)

	claims := jwt.MapClaims{
		"exp":    jwt.NewNumericDate(exp),
		"userID": user.ID.String(),
	}

	token, err := GenerateToken(claims, bus.tokenCfg.Secret, bus.tokenCfg.Method)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (bus Bussiness) Authenticate(ctx context.Context, tokenString string) (AuthenticateResp, error) {
	token, err := ParseToken(tokenString, bus.tokenCfg.Secret)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrSignatureInvalid) {
			return AuthenticateResp{}, ErrInvalidToken
		}
		return AuthenticateResp{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return AuthenticateResp{}, fmt.Errorf("unexpected claims type %T", token.Claims)
	}

	rawUserID, ok := claims["userID"].(string)

	if !ok {
		return AuthenticateResp{}, ErrUserIDMissing
	}

	userID, err := uuid.Parse(rawUserID)

	if err != nil {
		return AuthenticateResp{}, ErrUserIDMalformed
	}

	_, err = bus.userBus.GetUserByID(ctx, userID)

	if err != nil {
		return AuthenticateResp{}, ErrInvalidCredentials
	}

	return AuthenticateResp{
		UserID: userID,
	}, nil
}
