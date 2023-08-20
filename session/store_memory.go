package session

import (
	"github.com/shura1014/common/container/concurrent"
	"time"
)

// MemStore 内存中存储session，只适用于单节点应用
type MemStore struct {
	*concurrent.Map[string, *Session]
	ttl time.Duration
}

// SaveSession 存取session
func (store *MemStore) SaveSession(sessionId string, session *Session) {
	store.Put(sessionId, session)
}

// GetSession 获取session，如果已经过期那么返回一个空
func (store *MemStore) GetSession(sessionId string) *Session {
	value := store.Get(sessionId)
	if value == nil || value.IsExpire() {
		return nil
	}
	return value
}

// SessionExist 会话是否存在
func (store *MemStore) SessionExist(sessionId string) bool {
	value := store.Get(sessionId)
	if value.IsExpire() {
		return false
	}
	return true
}

// RemoveSession 移除会话
func (store *MemStore) RemoveSession(sessionId ...string) {
	store.Remove(sessionId...)
}

// ActiveCount 统计当前有效的session数量 ，并且同时会清理过期的session
func (store *MemStore) ActiveCount() int {
	count := 0
	evictKeys := make([]string, 0)
	store.Iterator(func(key string, session *Session) (ok bool) {
		if session.IsExpire() {
			evictKeys = append(evictKeys, key)
			// 清除过期
			return
		}
		count++
		return
	})

	store.RemoveSession(evictKeys...)
	return count
}

// CleanExpire 清理过期session
func (store *MemStore) CleanExpire() {
	store.ActiveCount()
}
