package main

import "sync"

type Storage struct {
	mu sync.RWMutex
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Set(key string, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val
 }

 func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
 }

 func (s *Storage) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
 }