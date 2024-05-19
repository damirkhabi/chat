package app

import (
	"context"
	"log"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/arifullov/chat-server/internal/clients/db"
	"github.com/arifullov/chat-server/internal/clients/db/pg"
	"github.com/arifullov/chat-server/internal/clients/db/transaction"
	"github.com/arifullov/chat-server/pkg/access_v1"

	"github.com/arifullov/chat-server/internal/api/chat"
	"github.com/arifullov/chat-server/internal/clients/grpc/auth"
	"github.com/arifullov/chat-server/internal/closer"
	"github.com/arifullov/chat-server/internal/config"
	"github.com/arifullov/chat-server/internal/repository"
	chatRepository "github.com/arifullov/chat-server/internal/repository/chat"
	"github.com/arifullov/chat-server/internal/service"
	chatService "github.com/arifullov/chat-server/internal/service/chat"
)

type serviceProvider struct {
	pgConfig         config.PGConfig
	grpcConfig       config.GRPCConfig
	authClientConfig config.AuthClientConfig
	jaegerConfig     config.JaegerConfig

	authClient auth.Client

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

func (s *serviceProvider) AuthClientConfig() config.AuthClientConfig {
	if s.authClientConfig == nil {
		authClientConfig, err := config.NewAuthClientConfig()
		if err != nil {
			log.Fatalf("failed to get auth client config: %s", err.Error())
		}
		s.authClientConfig = authClientConfig
	}
	return s.authClientConfig
}

func (s *serviceProvider) JaegerConfig() config.JaegerConfig {
	if s.jaegerConfig == nil {
		cfg, err := config.NewJaegerConfig()
		if err != nil {
			log.Fatalf("failed to get jaeger config: %s", err.Error())
		}
		s.jaegerConfig = cfg
	}
	return s.jaegerConfig
}

func (s *serviceProvider) AuthClient(ctx context.Context) auth.Client {
	if s.authClient == nil {
		creds := grpc.WithTransportCredentials(insecure.NewCredentials())

		conn, err := grpc.DialContext(
			ctx,
			s.AuthClientConfig().Address(),
			creds,
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)
		if err != nil {
			log.Fatalf("failed to connect to grpc %s: %s", s.AuthClientConfig().Address(), err.Error())
		}
		closer.Add(conn.Close)

		client := access_v1.NewAccessV1Client(conn)
		s.authClient = auth.NewClient(client)
	}
	return s.authClient
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
