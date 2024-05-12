package chat

import (
	"context"

	"github.com/arifullov/chat-server/internal/model"
)

func (s *serv) SendMessage(ctx context.Context, message *model.Message) error {
	err := s.chatRepository.CreateMessage(ctx, message)
	if err != nil {
		return err
	}
	return nil
}
