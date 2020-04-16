package accesslog

import (
	"sync"
)

type syncIPs struct {
	sync.RWMutex
	items map[int64][]string
}

func (s *syncIPs) all() map[int64][]string {
	s.RLock()
	defer s.RUnlock()
	return s.items
}

func (s *syncIPs) allById(id int64) []string {
	s.RLock()
	defer s.RUnlock()
	return s.items[id]
}

func (s *syncIPs) addIp(id int64, ip string) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.items[id]; !ok {
		s.items[id] = []string{ip}
		return
	}
	has := false
	for _, v := range s.items[id] {
		if v == ip {
			has = true
		}
	}
	if !has {
		s.items[id] = append(s.items[id], ip)
	}
}

func newSyncIPs() *syncIPs {
	return &syncIPs{items: make(map[int64][]string)}
}
