package config

import (
	"os"

	"github.com/pkg/errors"
)

const authAddressEnvName = "AUTH_ADDRESS"

type AuthClientConfig interface {
	Address() string
}

type authClientConfig struct {
	address string
}

func NewAuthClientConfig() (AuthClientConfig, error) {
	address := os.Getenv(authAddressEnvName)
	if address == "" {
		return nil, errors.New("address auth client not found")
	}
	return &authClientConfig{address: address}, nil
}

func (c authClientConfig) Address() string {
	return c.address
}
