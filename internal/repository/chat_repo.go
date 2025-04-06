package repository

import (
	"context"

	dbCtx "example.com/api/internal/repository/db"
)

type IChatRepo interface {
	CreateMessage(ctx context.Context, params dbCtx.CreateMessageParams) (dbCtx.Message, error)
	GetMessages(ctx context.Context, params dbCtx.GetMessagesParams) ([]dbCtx.GetMessagesRow, error)
}
