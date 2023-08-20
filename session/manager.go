package session

import (
	"github.com/shura1014/common/container/concurrent"
	"time"
)

const (
	MemoryStore    = "memory"
	RedisStore     = "redis"
	MemcachedStore = "memcached"
	defaultExpire  = 30 * time.Minute
	defaultStore   = MemoryStore
)

type Manager struct {
	store Store
	ttl   time.Duration
}

func New(ttl time.Duration, store string) *Manager {
	if ttl == 0 {
		ttl = defaultExpire
	}
	if store == "" {
		store = defaultStore
	}
	return &Manager{
		ttl:   ttl,
		store: NewStore(store, ttl),
	}
}

func Default() *Manager {
	return &Manager{
		ttl:   defaultExpire,
		store: NewStore(defaultStore, defaultExpire),
	}
}

func (manager *Manager) GetStore() Store {
	return manager.store
}

func (manager *Manager) NewSession() *Session {
	return &Session{
		id:      NewSessionId(),
		data:    concurrent.NewMap[string, any](),
		manager: manager,
		new:     true,
	}
}
