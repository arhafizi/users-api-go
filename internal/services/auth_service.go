package services

import (
	"context"

	dto "example.com/api/internal/contracts"
	dbCtx "example.com/api/internal/repository/db"
	"github.com/golang-jwt/jwt/v5"
)

type IAuthService interface {
	GenerateAccessToken(userID string) (string, error)

	GenerateRefreshToken(userID string) (string, error)

	ValidateToken(tokenString string) (jwt.MapClaims, error)

	Authenticate(ctx context.Context, email, password string) (*dbCtx.User, error)

	RegisterUser(ctx context.Context, args dto.Register) (*dto.UserResponse, string, string, error)
}
