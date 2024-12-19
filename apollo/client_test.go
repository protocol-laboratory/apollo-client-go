package apollo

import (
	"github.com/libgox/addr"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient(&Config{
		AppID:          "SampleApp",
		Cluster:        "default",
		NamespaceNames: []string{"application", "application2"},
		Address: addr.Address{
			Host: "localhost",
			Port: 8080,
		},
		Secret:    "",
		TLSConfig: nil,
		Logger:    nil,
	})
	if err == nil {
		value := c.GetStringValue("application", "timeout")
		value2 := c.GetStringValue("application2", "timeout")
		c.SubscribeEvent(&ClientTest{})
		t.Log(value, ",", value2)
	}
	time.Sleep(100 * time.Second)
}

type ClientTest struct{}

func (c *ClientTest) OnChange(event *ChangeEvent) {
}
