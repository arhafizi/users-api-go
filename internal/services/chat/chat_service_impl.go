package chat

import (
	"context"

	"example.com/api/internal/repository"
	dbCtx "example.com/api/internal/repository/db"
	"example.com/api/pkg/logging"
)

type ChatService struct {
	repo   repository.IRepositoryManager
	logger logging.ILogger
}

func NewChatService(r repository.IRepositoryManager, l logging.ILogger) *ChatService {
	return &ChatService{
		repo:   r,
		logger: l,
	}
}

func (s *ChatService) SaveMessage(ctx context.Context, senderID int32, content string) error {
	_, err := s.repo.Chat().CreateMessage(ctx, dbCtx.CreateMessageParams{
		SenderID: senderID,
		Content:  content,
	})
	return err
}

func (s *ChatService) GetMessages(ctx context.Context, limit, offset int32) ([]dbCtx.GetMessagesRow, error) {
	return s.repo.Chat().GetMessages(ctx, dbCtx.GetMessagesParams{
		Limit:  limit,
		Offset: offset,
	})
}
