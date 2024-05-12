package chat

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/arifullov/chat-server/internal/converter"
	desc "github.com/arifullov/chat-server/pkg/chat_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.chatService.Create(ctx, converter.ToChatFromDesc(req))
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to create chat")
	}
	log.Printf("inserted chat with id: %d", id)
	return &desc.CreateResponse{
		Id: id,
	}, nil
}
