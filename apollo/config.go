package apollo

import (
	"crypto/tls"
	"fmt"
	"github.com/libgox/addr"
	"net/url"
)

type Config struct {
	AppID          string
	Cluster        string
	NamespaceNames []string
	Address        addr.Address
	Secret         string
	// TlsConfig configuration information for tls.
	TLSConfig *tls.Config
	Logger    Logger
}

func (c *Config) GetNotifyURLSuffix(notifications string) string {
	return fmt.Sprintf("%s/notifications/v2?appId=%s&cluster=%s&notifications=%s",
		c.GetUrlPrefix(),
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.Cluster),
		url.QueryEscape(notifications))
}

func (c *Config) GetSyncURI(namespace string) string {
	return fmt.Sprintf("%s/configs/%s/%s/%s?releaseKey=&ip=%s",
		c.GetUrlPrefix(),
		url.QueryEscape(c.AppID),
		url.QueryEscape(c.Cluster),
		url.QueryEscape(namespace),
		GetLocalIP())
}

func (c *Config) GetUrlPrefix() string {
	var urlPrefix string
	if c.TLSConfig != nil {
		urlPrefix = "https://" + c.Address.Addr()
	} else {
		urlPrefix = "http://" + c.Address.Addr()
	}
	return urlPrefix
}
