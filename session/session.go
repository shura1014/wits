package session

import (
	"github.com/shura1014/common/container/concurrent"
	"github.com/shura1014/common/random/id"
	"github.com/shura1014/common/utils/timeutil"
	"time"
)

func NewSessionId() string {
	return id.Str()
}

type Session struct {
	id             string
	dirty          bool
	data           *concurrent.Map[string, any]
	new            bool
	manager        *Manager
	lastAccessTime time.Time
}

func (s *Session) GetId() string {
	s.lastAccessTime = time.Now()
	return s.id
}

func (s *Session) IsNew() bool {
	return s.new
}

func (s *Session) IsDirty() bool {
	return s.dirty
}

func (s *Session) Get(name string) any {
	s.lastAccessTime = time.Now()
	return s.data.Get(name)
}

func (s *Session) GetAll() map[string]any {
	s.lastAccessTime = time.Now()
	return s.data.GetAll()
}

func (s *Session) Set(name string, value any) {
	s.lastAccessTime = time.Now()
	s.data.Put(name, value)
	s.dirty = true
}

func (s *Session) SetMap(data map[string]any) {
	s.lastAccessTime = time.Now()
	s.data.PutAll(data)
	s.dirty = true
}

func (s *Session) Remove(name string) {
	s.lastAccessTime = time.Now()
	s.data.Remove(name)
	s.dirty = true
}

func (s *Session) ClearState() {
	s.dirty = false
	s.new = false
}

func (s *Session) Save() {
	s.lastAccessTime = time.Now()
	s.manager.store.SaveSession(s.id, s)
	s.ClearState()
}

// defaultExpire 默认的过期时间
// appoint 指定的过期时间
func computeExpire(defaultExpire int64, appoint ...int64) int64 {
	expire := defaultExpire
	if len(appoint) > 0 {
		expire = appoint[0]
	}

	return timeutil.MilliSeconds() + expire

}

func (s *Session) IsExpire() bool {
	return s.lastAccessTime.Add(s.manager.ttl).Before(time.Now())
}
