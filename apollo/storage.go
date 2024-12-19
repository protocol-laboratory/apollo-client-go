package apollo

import (
	"sync"
)

// storage namespace cache
type storage struct {
	caches sync.Map
}

func newStorage(namespaceNames []string) *storage {
	s := &storage{
		caches: sync.Map{},
	}
	for _, namespace := range namespaceNames {
		s.caches.Store(namespace, &cache{
			data: sync.Map{},
		})
	}
	return s
}

func (s *storage) loadCache(namespace string) *cache {
	if value, ok := s.caches.Load(namespace); ok {
		return value.(*cache)
	}
	c := &cache{
		data: sync.Map{},
	}
	s.caches.Store(namespace, c)
	return c
}

// apolloConfiguration query config result
type apolloConfiguration struct {
	NamespaceName  string            `json:"namespaceName"`
	Configurations map[string]string `json:"configurations"`
	ReleaseKey     string            `json:"releaseKey"`
}

// cache apollo namespace configuration cache
type cache struct {
	data sync.Map
}
