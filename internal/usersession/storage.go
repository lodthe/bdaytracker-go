package usersession

import (
	"sync"
)

type storage struct {
	sessionLockers map[int]sync.Locker
	lock           sync.Locker
}

func newStorage() *storage {
	return &storage{
		sessionLockers: map[int]sync.Locker{},
		lock:           &sync.Mutex{},
	}
}

func (s *storage) acquireLock(telegramID int) sync.Locker {
	s.lock.Lock()
	defer s.lock.Unlock()

	sessionLocker, exists := s.sessionLockers[telegramID]
	if !exists {
		sessionLocker = &sync.Mutex{}
		s.sessionLockers[telegramID] = sessionLocker
	}

	return sessionLocker
}
