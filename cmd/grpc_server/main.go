package main

import (
	"context"
	"fmt"
	"log"
	"net"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/arifullov/chat-server/pkg/chat_v1"
)

const grpcPort = 50053
const dbDSN = "host=localhost port=5002 dbname=chat_db user=chat password=secret_pass sslmode=disable"

type server struct {
	desc.UnimplementedChatV1Server
	dbPool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("create chat: %v", req)
	builderInsert := sq.Insert("chats").
		PlaceholderFormat(sq.Dollar).
		Columns("usernames").
		Values(req.Usernames).
		Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to build query")
	}

	var chatID int64
	err = s.dbPool.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		log.Printf("failed to insert chat: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to insert chat")
	}
	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	builderDelete := sq.Delete("chats").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.Id})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to delete chat")
	}

	res, err := s.dbPool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to delete chat: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to delete chat")
	}
	if res.RowsAffected() == 0 {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}
	log.Printf("delete chat: %v", req)
	return nil, nil
}
func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("send message: %v", req)
	builderInsert := sq.Insert("messages").
		PlaceholderFormat(sq.Dollar).
		Columns("\"from\"", "text", "chat_id", "timestamp").
		Values(req.From, req.Text, req.ChatID, req.Timestamp.AsTime()).
		Suffix("RETURNING id")
	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to build query")
	}

	var chatID int64
	err = s.dbPool.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		log.Printf("failed to insert messages: %v", err)
		return nil, status.Errorf(codes.Unavailable, "failed to insert messages")
	}
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed db connect: %v", err)
	}
	if err = dbPool.Ping(ctx); err != nil {
		log.Fatalf("failed db connect: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{dbPool: dbPool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
