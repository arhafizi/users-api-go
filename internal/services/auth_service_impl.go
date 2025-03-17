package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"example.com/api/config"
	dbCtx "example.com/api/internal/repository/db"
	"example.com/api/internal/services/hashing"
	"example.com/api/pkg/logging"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	logger      logging.ILogger
	jwtConf     config.JWTConfig
	hashService hashing.IHashService
	userService IUserService
}

func NewAuthService(jwtConf config.JWTConfig, hasher hashing.IHashService, userService IUserService, logger logging.ILogger) *AuthService {
	return &AuthService{
		logger:      logger,
		jwtConf:     jwtConf,
		hashService: hasher,
		userService: userService,
	}
}

func (s *AuthService) GenerateAccessToken(userID string) (string, error) {
	expTime := time.Now().Add(s.jwtConf.AccessTokenExpireDuration * time.Minute).Unix()

	claims := jwt.MapClaims{
		"user_id":    userID,
		"exp":        expTime,
		"iat":        time.Now().Unix(),
		"token_type": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConf.Secret))
}

func (s *AuthService) GenerateRefreshToken(userID string) (string, error) {
	expTime := time.Now().Add(s.jwtConf.RefreshTokenExpireDuration * time.Minute).Unix()

	claims := jwt.MapClaims{
		"user_id":    userID,
		"exp":        expTime,
		"iat":        time.Now().Unix(),
		"token_type": "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtConf.Secret))
}

func (s *AuthService) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.jwtConf.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *AuthService) Authenticate(ctx context.Context, email, password string) (*dbCtx.User, error) {
	user, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to fetch user")
	}

	if err = s.hashService.Compare(user.PasswordHash, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, name, email, password string) (*dbCtx.User, string, string, error) {
	hashedPassword, err := s.hashService.Hash(password)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to hash password: %w", err)
	}
	createParams := dbCtx.CreateUserParams{
		Username:     name,
		Email:        email,
		PasswordHash: hashedPassword,
	}
	user, err := s.userService.Create(ctx, createParams)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create user: %w", err)
	}

	accessToken, err := s.GenerateAccessToken(fmt.Sprintf("%d", user.ID))
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.GenerateRefreshToken(fmt.Sprintf("%d", user.ID))
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return user, accessToken, refreshToken, nil
}
