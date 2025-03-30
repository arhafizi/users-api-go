package storage

import (
	"context"
	"time"
)

type ITokenStorage interface {
	Store(ctx context.Context, userID string, tokenID string, exp time.Duration) error
	Validate(ctx context.Context, userID string, tokenID string) error
	Invalidate(ctx context.Context, userID string) error
}
