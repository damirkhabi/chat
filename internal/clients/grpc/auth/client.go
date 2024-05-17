package auth

import (
	"context"

	"github.com/arifullov/chat-server/pkg/access_v1"
)

var _ = (*client)(nil)

type Client interface {
	Check(ctx context.Context, endpoint string) (bool, error)
}

type client struct {
	accessClient access_v1.AccessV1Client
}

func NewClient(accessClient access_v1.AccessV1Client) Client {
	return &client{accessClient: accessClient}
}

func (c client) Check(ctx context.Context, endpoint string) (bool, error) {
	if _, err := c.accessClient.Check(ctx, &access_v1.CheckRequest{
		EndpointAddress: endpoint,
	}); err != nil {
		return false, err
	}

	return true, nil
}
