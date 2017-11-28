package manager

import (
	"sync"
	"wxbot2/wxbot"
)

var (
	GlobalSessionManager = &SessionManager{
		sm: make(map[string]*wxbot.Wxbot),
	}
)

type SessionManager struct {
	sm   map[string]*wxbot.Wxbot
	lock sync.RWMutex
}

func (s *SessionManager) Set(uuid string,session *wxbot.Wxbot) string {
	// generate uuid

	s.lock.Lock()
	s.sm[uuid] = session
	s.lock.Unlock()

	return uuid
}

func (s *SessionManager) Get(uuid string) *wxbot.Wxbot {
	s.lock.RLock()
	session := s.sm[uuid]
	s.lock.RUnlock()
	return session
}


