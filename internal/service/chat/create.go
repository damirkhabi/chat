package chat

import (
	"context"

	"github.com/arifullov/chat-server/internal/model"
)

func (s *serv) Create(ctx context.Context, chat *model.Chat) (int64, error) {
	chadID, err := s.chatRepository.Create(ctx, chat)
	if err != nil {
		return 0, err
	}
	return chadID, nil
}
