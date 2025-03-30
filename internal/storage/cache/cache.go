package cache

import (
	"context"
	"time"
)

type ICacheService interface {
	Get(ctx context.Context, key string, dest any) (bool, error)
	Set(ctx context.Context, key string, value any, exp time.Duration) error
	Delete(ctx context.Context, key string) error
}
