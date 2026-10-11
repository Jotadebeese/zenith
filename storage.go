package main

import (
	"sync"
	"time"
)

type Item struct {
	val string
	expiresAt time.Time
	vector []float32
}

type Storage struct {
	mu sync.RWMutex
	data map[string]Item
}

type VectorMatch struct {
	Key string
	Val string
	Score float32
}

func NewStorage() *Storage {
	s := &Storage{
		data: make(map[string]Item),
	}
	go s.startCleaner(100 * time.Millisecond)
	return s
}
func (s *Storage) Set(key string, val string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	s.data[key] = Item{
		val: val,
		expiresAt: expiresAt,
	}
 }

 func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	item, ok := s.data[key]
	s.mu.RUnlock()

	if !ok {
		return "", false
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		s.Del(key)
		return "", false
	}

	return item.val, ok
 }

 func (s *Storage) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
 }

 func (s *Storage) startCleaner(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, item := range s.data {
			if !item.expiresAt.IsZero() && now.After(item.expiresAt) {
				delete(s.data, key)
			}
		}
		s.mu.Unlock()
	}
 }


  func (s *Storage) SetVector(key string, val string, vector []float32, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	s.data[key] = Item{
		val: val,
		expiresAt: expiresAt,
		vector: vector,
	}
  }

  func (s *Storage) QueryVector(target []float32, threshold float32) (*VectorMatch, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	var bestMatch *VectorMatch
	bestScore := float32(-1.0)

	for key, item := range s.data {
		if !item.expiresAt.IsZero() && now.After(item.expiresAt) {
			continue
		}

		if len(item.vector) == 0 {
			continue
		}

		score, err := CosineSimilarity(target, item.vector)
		if err == nil && score >= threshold && score > bestScore {
			bestScore = score
			bestMatch = &VectorMatch{
				Key: key,
				Val: item.val,
				Score: score,
			}
		}
	}

	if bestMatch == nil {
		return nil, false
	}
	return bestMatch, true
  }