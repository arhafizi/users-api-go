package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisTokenStorage struct {
	client    *redis.Client
	keyPrefix string
}

func NewRedisTokenStorage(c *redis.Client, p string) *RedisTokenStorage {
	return &RedisTokenStorage{
		client:    c,
		keyPrefix: p,
	}
}

func (s *RedisTokenStorage) key(userID string) string {
	return fmt.Sprintf("user-%s", userID)
}

func (s *RedisTokenStorage) Store(ctx context.Context, userID, tokenID string, exp time.Duration) error {
	return s.client.Set(ctx, s.key(userID), tokenID, exp).Err()
}

func (s *RedisTokenStorage) Validate(ctx context.Context, userID, tokenID string) error {
	storedID, err := s.client.Get(ctx, s.key(userID)).Result()
	if err == redis.Nil {
		return fmt.Errorf("token not found")
	}
	if err != nil {
		return fmt.Errorf("storage error: %w", err)
	}
	if storedID != tokenID {
		return fmt.Errorf("token mismatch")
	}
	return nil
}

func (s *RedisTokenStorage) Invalidate(ctx context.Context, userID string) error {
	return s.client.Del(ctx, s.key(userID)).Err()
}
