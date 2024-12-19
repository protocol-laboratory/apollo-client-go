package apollo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultLongPollInterval = time.Second * 2

type notifyHandler func(n *notification) error

type requester struct {
	client  *http.Client
	retries int
}

func (r *requester) do(ctx context.Context, uri string, retries int) ([]byte, error) {
	resp, err := r.client.Get(uri)
	if err != nil {
		if retries > 0 {
			return r.do(ctx, uri, retries-1)
		}
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("failed to close response body: %v", err)
			return
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		return io.ReadAll(resp.Body)
	}

	return nil, fmt.Errorf("apollo return http resp code %d", resp.StatusCode)
}

type longPoll struct {
	config       *Config
	interval     time.Duration
	handler      notifyHandler
	requester    *requester
	notification *notificationsMgr
}

func newLongPoll(config *Config, handler notifyHandler) *longPoll {
	p := &longPoll{
		config:   config,
		interval: defaultLongPollInterval,
		handler:  handler,
		requester: &requester{
			client: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: config.TLSConfig,
				},
			},
			retries: 3,
		},
		notification: newNotificationManager(config.NamespaceNames),
	}
	return p
}

func (p *longPoll) start(ctx context.Context) {
	child, cancel := context.WithCancel(ctx)
	defer cancel()

	timer := time.NewTimer(p.interval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			err := p.fetch(child)
			log.Errorf("fetch config err: %v", err)
			timer.Reset(p.interval)
		case <-child.Done():
			return
		}
	}
}

func (p *longPoll) fetch(ctx context.Context) error {
	url := p.config.GetNotifyURLSuffix(p.notification.String())
	result, err := p.requester.do(ctx, url, p.requester.retries)
	if err != nil {
		return err
	}
	if len(result) == 0 {
		log.Warn("apollo get notify result empty")
		return nil
	}
	var n []*notification
	if err := json.Unmarshal(result, &n); err != nil {
		return err
	}
	for _, v := range n {
		if err := p.handler(v); err != nil {
			return err
		}
		p.notification.Store(v.NamespaceName, v.NotificationID)
	}
	return nil
}
