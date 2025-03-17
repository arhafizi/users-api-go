package services

import (
	"context"

	dbCtx "example.com/api/internal/repository/db"
	"github.com/golang-jwt/jwt/v5"
)

type IAuthService interface {
	GenerateAccessToken(userID string) (string, error)

	GenerateRefreshToken(userID string) (string, error)

	ValidateToken(tokenString string) (jwt.MapClaims, error)

	Authenticate(ctx context.Context, email, password string) (*dbCtx.User, error)

	RegisterUser(ctx context.Context, name, email, password string) (*dbCtx.User, string, string, error)
}
