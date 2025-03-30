package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"example.com/api/config"
	dto "example.com/api/internal/contracts"
	dbCtx "example.com/api/internal/repository/db"
	"example.com/api/internal/services/hashing"
	"example.com/api/internal/storage"
	"example.com/api/pkg/logging"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	logger       logging.ILogger
	jwtConf      config.JWTConfig
	hashService  hashing.IHashService
	userService  IUserService
	tokenStorage storage.ITokenStorage
}

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

func NewAuthService(
	jwtConf config.JWTConfig, hasher hashing.IHashService,
	userSvc IUserService, logger logging.ILogger,
	storage storage.ITokenStorage,
) *AuthService {
	return &AuthService{
		logger:       logger,
		jwtConf:      jwtConf,
		hashService:  hasher,
		userService:  userSvc,
		tokenStorage: storage,
	}
}

func (s *AuthService) GenerateAccessToken(userID string) (string, error) {
	expTime := time.Now().Add(s.jwtConf.AccessTokenExpireDuration * time.Minute).Unix()

	claims := jwt.MapClaims{
		"sub":        userID,
		"exp":        expTime,
		"iat":        time.Now().Unix(),
		"token_type": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConf.Secret))
}

func (s *AuthService) GenerateRefreshToken(userID string) (string, error) {
	tokenID := uuid.New().String()
	expTime := time.Now().Add(s.jwtConf.RefreshTokenExpireDuration * time.Minute).Unix()
	remainingTime := expTime - time.Now().Unix()
	err := s.tokenStorage.Store(context.Background(), userID, tokenID, time.Duration(remainingTime)*time.Second)
	if err != nil {
		return "", errors.New("failed to store refresh token")
	}
	claims := jwt.MapClaims{
		"sub":        userID,
		"exp":        expTime,
		"iat":        time.Now().Unix(),
		"token_type": "refresh",
		"id":         tokenID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.jwtConf.RefreshSecret))
}

func (s *AuthService) Authenticate(ctx context.Context, email, password string) (*dbCtx.User, error) {
	// email = strings.ToLower(email) //
	user, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		// s.logger.Error("Failed to fetch user by email: %v", err) //
		return nil, errors.New("invalid credentials, fetch")
	}

	if err := s.hashService.Compare(user.PasswordHash, password); err != nil {
		// s.logger.Error("Password comparison failed: %v", err)
		return nil, errors.New(err.Error())
	}

	return user, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, args dto.Register) (*dto.UserResponse, string, string, error) {
	createParams := dto.CreateUserReq{
		Username: args.Name,
		FullName: "",
		Email:    args.Email,
		Password: args.Password,
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

func (s *AuthService) RotateTokens(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.ValidateRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	userID := claims["sub"].(string)

	if err := s.tokenStorage.Invalidate(ctx, userID); err != nil {
		return "", "", err
	}

	accessToken, err := s.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) ValidateToken(tokenString string, expectedType string) (jwt.MapClaims, error) {
	var secretKey []byte
	switch expectedType {
	case tokenTypeAccess:
		secretKey = []byte(s.jwtConf.Secret)
	case tokenTypeRefresh:
		secretKey = []byte(s.jwtConf.RefreshSecret)
	default:
		return nil, errors.New("invalid token type specified")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func (s *AuthService) ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	return s.ValidateToken(tokenString, tokenTypeAccess)
}

func (s *AuthService) ValidateRefreshToken(ctx context.Context, tokenString string) (jwt.MapClaims, error) {
	claims, err := s.ValidateToken(tokenString, tokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	tokenID, ok := claims["id"].(string)
	if !ok {
		return nil, errors.New("invalid token ID")
	}

	if err := s.tokenStorage.Validate(ctx, userID, tokenID); err != nil {
		return nil, errors.New("invalid or revoked refresh token: " + err.Error())
	}

	return claims, nil
}

func (s *AuthService) InvalidateRefreshToken(ctx context.Context, userID string) error {
	return s.tokenStorage.Invalidate(ctx, userID)
}
