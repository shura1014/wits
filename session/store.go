package session

import (
	"github.com/shura1014/common/container/concurrent"
	"time"
)

type Store interface {
	SaveSession(sessionId string, session *Session)
	GetSession(sessionId string) *Session
	SessionExist(sessionId string) bool
	RemoveSession(sessionId ...string)
	ActiveCount() int
	CleanExpire()
}

func NewStore(store string, ttl time.Duration) Store {
	switch store {
	case "memory":
		return &MemStore{Map: concurrent.NewMap[string, *Session](), ttl: ttl}
	default:
		return &MemStore{Map: concurrent.NewMap[string, *Session](), ttl: ttl}
	}
}
