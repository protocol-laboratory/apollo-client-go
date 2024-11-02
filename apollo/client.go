package apollo

import (
	"crypto/tls"

	"github.com/libgox/addr"
)

type Client interface {
}

type Config struct {
	Address addr.Address
	// TlsConfig configuration information for tls.
	TlsConfig *tls.Config
}

func NewClient(config *Config) (Client, error) {
	return newClient(config)
}
