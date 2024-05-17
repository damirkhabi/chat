package chat

import (
	"github.com/arifullov/chat-server/internal/clients/db"
	"github.com/arifullov/chat-server/internal/repository"
	"github.com/arifullov/chat-server/internal/service"
)

type serv struct {
	chatRepository repository.ChatRepository
	txManager      db.TxManager
}

func NewChatService(
	chatRepository repository.ChatRepository,
	txManager db.TxManager,
) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
		txManager:      txManager,
	}
}
