package repository

import (
	"context"

	"github.com/arifullov/chat-server/internal/model"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *model.Chat) (int64, error)
	Delete(ctx context.Context, id int64) error
	CreateMessage(ctx context.Context, message *model.Message) error
}
