package repository

import (
	"context"

	dbCtx "example.com/api/internal/repository/db"
)

type ChatRepository struct {
	q *dbCtx.Queries
}

func NewChatRepo(q dbCtx.DBTX) *ChatRepository {
	return &ChatRepository{q: dbCtx.New(q)}
}

func (r *ChatRepository) CreateMessage(ctx context.Context, params dbCtx.CreateMessageParams) (dbCtx.Message, error) {
	return r.q.CreateMessage(ctx, params)
}

func (r *ChatRepository) GetMessages(ctx context.Context, params dbCtx.GetMessagesParams) ([]dbCtx.GetMessagesRow, error) {
	return r.q.GetMessages(ctx, params)
}
