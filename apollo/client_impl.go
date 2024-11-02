package apollo

type client struct {
	urlPrefix string
}

func newClient(config *Config) (*client, error) {
	c := &client{}
	if config.TlsConfig != nil {
		c.urlPrefix = "https://" + config.Address.Addr()
	} else {
		c.urlPrefix = "http://" + config.Address.Addr()
	}
	return c, nil
}
