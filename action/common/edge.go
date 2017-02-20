package common

import (
	"sync"
)

type EdgeStore struct {
	State map[string]int
	Lock sync.Mutex
}

func (s *EdgeStore) Delta(ID string, state int) bool {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	if s.State == nil {
		s.State = map[string]int{}
	}
	c, ok := s.State[ID]
	if !ok || c != state {
		s.State[ID] = state
		return true
	}
	return false
}
