package apollo

import (
	"context"
	"crypto/tls"

	"github.com/libgox/addr"
	"golang.org/x/exp/slog"
)

type Client interface {
}

type innerClient struct {
	ctx context.Context
}

type Config struct {
	Address addr.Address
	// TlsConfig configuration information for tls.
	TlsConfig *tls.Config
	// Logger structured logger for logging operations
	Logger *slog.Logger
}

func NewClient(config *Config) (Client, error) {
	if config.Logger != nil {
		config.Logger = slog.Default()
	}
	ctx := context.Background()
	c := &innerClient{
		ctx: ctx,
	}
	return c, nil
}
