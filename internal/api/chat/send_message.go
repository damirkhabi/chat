package chat

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/arifullov/chat-server/internal/converter"
	desc "github.com/arifullov/chat-server/pkg/chat_v1"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	err := i.chatService.SendMessage(ctx, converter.ToMessageFromDesc(req))
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to send message")
	}
	log.Printf("send message: %v", req)
	return nil, nil
}
