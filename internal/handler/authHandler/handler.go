package authhandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	userbus "github.com/Javlon721/booklib/internal/bus/userBus"
	"github.com/Javlon721/booklib/internal/handler/errs"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidAuthParams  = errors.New("invalid auth params")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserIDMissing      = errors.New("userID missing from token hander")
	ErrUserIDMalformed    = errors.New("userID in token header is malformed")
)

type UserBus interface {
	GetUserByID(context.Context, uuid.UUID) (userbus.User, error)
	GetUserByEmail(context.Context, string) (userbus.User, error)
}

type Handler struct {
	userBus  UserBus
	logger   *slog.Logger
	tokenCfg TokenConfig
}

func (h Handler) Login(c fiber.Ctx) error {
	var payload AuthParams

	if err := c.Bind().Body(&payload); err != nil {
		return errs.New(errs.InvalidArgument, ErrInvalidAuthParams)
	}

	user, err := h.userBus.GetUserByEmail(c.Context(), payload.Email)

	if err != nil {
		if errors.Is(err, userbus.ErrUserNotFound) {
			return errs.New(errs.NotFound, userbus.ErrUserNotFound)
		}

		return errs.New(errs.Internal, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password)); err != nil {
		return errs.New(errs.PermissionDenied, ErrInvalidCredentials)
	}

	exp := time.Now().UTC().Add(h.tokenCfg.TokenExpiresAt)

	claims := jwt.MapClaims{
		"exp":    jwt.NewNumericDate(exp),
		"userID": user.ID.String(),
	}

	token, err := GenerateToken(claims, h.tokenCfg.Secret, h.tokenCfg.Method)

	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return c.SendString(token)
}

func (h Handler) Authenticate(ctx context.Context, tokenString string) (AuthenticateResp, error) {
	token, err := ParseToken(tokenString, h.tokenCfg.Secret)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrSignatureInvalid) {
			return AuthenticateResp{}, errs.New(errs.Unauthenticated, err)
		}
		return AuthenticateResp{}, errs.New(errs.Internal, err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return AuthenticateResp{}, errs.New(errs.Internal, fmt.Errorf("unexpected claims type"))
	}

	rawUserID, ok := claims["userID"].(string)

	if !ok {
		return AuthenticateResp{}, errs.New(errs.Unauthenticated, ErrUserIDMissing)
	}

	userID, err := uuid.Parse(rawUserID)

	if err != nil {
		return AuthenticateResp{}, errs.New(errs.Unauthenticated, ErrUserIDMalformed)
	}

	_, err = h.userBus.GetUserByID(ctx, userID)

	if err != nil {
		return AuthenticateResp{}, errs.New(errs.Unauthenticated, ErrInvalidCredentials)
	}

	return AuthenticateResp{
		UserID: userID,
	}, nil
}

func NewHandler(userBus UserBus, logger *slog.Logger, cfg TokenConfig) *Handler {
	return &Handler{
		userBus:  userBus,
		logger:   logger,
		tokenCfg: cfg,
	}
}
