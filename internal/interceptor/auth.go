package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/arifullov/chat-server/internal/clients/grpc/auth"
)

type authInterceptor struct {
	authClient auth.Client
}

func NewAuthInterceptor(authClient auth.Client) *authInterceptor {
	return &authInterceptor{authClient: authClient}
}

func (i *authInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		_, err := i.authClient.Check(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
