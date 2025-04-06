package chat

import (
	"context"

	dbCtx "example.com/api/internal/repository/db"
)

type IChatService interface {
	SaveMessage(ctx context.Context, senderID int32, content string) error
	GetMessages(ctx context.Context, limit, offset int32) ([]dbCtx.GetMessagesRow, error)
}
