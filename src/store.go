package main

import "sync"

// URLStore
type URLStore struct {
	urls map[string]string
	m    sync.RWMutex
}

func (s *URLStore) Get(key string) (val string, ok bool) {
	s.m.RLock()
	defer s.m.RUnlock()
	val, ok = s.urls[key]
	return
}

func (s *URLStore) Set(key, url string) bool {
	s.m.Lock()
	defer s.m.Unlock()
	_, present := s.urls[key]
	if present {
		return false
	}
	s.urls[key] = url
	return true
}

func (s *URLStore) Count() int {
	s.m.RLock()
	defer s.m.RUnlock()
	return len(s.urls)
}

func (s *URLStore) Push(url string) string {
	for {
		key := genKey(s.Count())
		if s.Set(key, url) {
			return key
		}
	}
}

func NewURLStore() *URLStore {
	return &URLStore{urls: make(map[string]string)}
}
