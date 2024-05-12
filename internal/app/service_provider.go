package app

import (
	"context"
	"log"

	"github.com/arifullov/chat-server/internal/api/chat"
	"github.com/arifullov/chat-server/internal/client/db"
	"github.com/arifullov/chat-server/internal/client/db/pg"
	"github.com/arifullov/chat-server/internal/client/db/transaction"
	"github.com/arifullov/chat-server/internal/closer"
	"github.com/arifullov/chat-server/internal/config"
	"github.com/arifullov/chat-server/internal/repository"
	"github.com/arifullov/chat-server/internal/service"

	chatRepository "github.com/arifullov/chat-server/internal/repository/chat"
	chatService "github.com/arifullov/chat-server/internal/service/chat"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	chatRepository repository.ChatRepository

	chatService service.ChatService

	userImpl *chat.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		pgConfig, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}
		s.pgConfig = pgConfig
	}
	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		grpcConfig, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = grpcConfig
	}
	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		if err = cl.DB().Ping(ctx); err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}
	return s.dbClient
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chatRepository.NewRepository(s.DBClient(ctx))
	}
	return s.chatRepository
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}
	return s.txManager
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewChatService(
			s.ChatRepository(ctx),
			s.TxManager(ctx),
		)
	}
	return s.chatService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *chat.Implementation {
	if s.userImpl == nil {
		s.userImpl = chat.NewImplementation(s.ChatService(ctx))
	}
	return s.userImpl
}
