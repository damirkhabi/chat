package chat

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/arifullov/chat-server/internal/model"
	desc "github.com/arifullov/chat-server/pkg/chat_v1"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := i.chatService.Delete(ctx, req.GetId())
	if errors.Is(err, model.ErrChatNotFound) {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to create chat")
	}
	log.Printf("deleted chat with id: %d", req.GetId())
	return nil, nil
}
