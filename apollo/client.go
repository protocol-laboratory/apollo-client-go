package apollo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	GetStringValue(namespace, key string) string

	SubscribeEvent(listener Listener)
}

type innerClient struct {
	ctx      context.Context
	config   *Config
	storage  *storage
	poller   *longPoll
	listener Listener
}

func NewClient(config *Config) (Client, error) {
	if config.Logger != nil {
		SetLogger(config.Logger)
	}
	ctx := context.Background()
	c := &innerClient{
		ctx:     ctx,
		config:  config,
		storage: newStorage(config.NamespaceNames),
	}
	c.poller = newLongPoll(config, c.updateHandle)

	// sync
	err := c.poller.fetch(c.ctx)
	if err != nil {
		return nil, err
	}

	// long poll
	go c.poller.start(c.ctx)

	return c, nil
}

func (i *innerClient) updateHandle(notification *notification) error {
	change, err := i.sync(notification)
	if err != nil {
		return err
	}
	if change == nil || len(change.Changes) == 0 {
		return fmt.Errorf("no changes to sync")
	}
	if i.listener != nil {
		i.listener.OnChange(change)
	}
	return nil
}

func (i *innerClient) sync(notification *notification) (*ChangeEvent, error) {
	log.Infof("sync namespace %s with remote config server", notification.NamespaceName)
	url := i.config.GetSyncURI(notification.NamespaceName)
	r := &requester{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: i.config.TLSConfig,
			},
		},
		retries: 3,
	}
	result, err := r.do(i.ctx, url, r.retries)
	if err != nil || len(result) == 0 {
		return nil, err
	}

	ac := &apolloConfiguration{}
	if err = json.Unmarshal(result, ac); err != nil {
		return nil, err
	}
	return i.updateCache(ac)
}

func (i *innerClient) updateCache(ac *apolloConfiguration) (*ChangeEvent, error) {
	var change = &ChangeEvent{
		Namespace: ac.NamespaceName,
		Changes:   make(map[string]*Change),
	}
	c := i.storage.loadCache(ac.NamespaceName)

	c.data.Range(func(k, v interface{}) bool {
		key := k.(string)
		value := v.(string)
		if _, ok := ac.Configurations[key]; !ok {
			c.data.Delete(key)
			change.Changes[key] = onDelete(key, value)
		}
		return true
	})

	for k, v := range ac.Configurations {
		old, ok := c.data.Load(k)
		if !ok {
			change.Changes[k] = onAdd(k, v)
			c.data.Store(k, v)
			continue
		}
		if old.(string) != v {
			change.Changes[k] = onModify(k, old.(string), v)
		}
		c.data.Store(k, v)
	}
	return change, nil
}

func (i *innerClient) SubscribeEvent(listener Listener) {
	i.listener = listener
}

func (i *innerClient) GetStringValue(namespace string, key string) string {
	v, ok := i.storage.loadCache(namespace).data.Load(key)
	if !ok {
		return ""
	}
	return v.(string)
}
