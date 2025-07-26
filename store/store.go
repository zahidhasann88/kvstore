package store

import (
	"sync"
	"time"
)

type Item struct {
	Value     string    `json:"value"`
	ExpiresAt time.Time `json:"expires_at"`
	HasTTL    bool      `json:"has_ttl"`
}

type Store struct {
	data        map[string]*Item
	mutex       sync.RWMutex
	expiration  *ExpirationManager
	persistence *PersistenceManager
}

func NewStore() *Store {
	s := &Store{
		data: make(map[string]*Item),
	}
	s.expiration = NewExpirationManager(s)
	s.persistence = NewPersistenceManager(s)
	return s
}

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	item := &Item{
		Value:  value,
		HasTTL: ttl > 0,
	}

	if ttl > 0 {
		item.ExpiresAt = time.Now().Add(ttl)
	}

	s.data[key] = item

	if ttl > 0 {
		s.expiration.SetExpiration(key, ttl)
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	item, exists := s.data[key]
	if !exists {
		return "", false
	}

	if item.HasTTL && time.Now().After(item.ExpiresAt) {
		delete(s.data, key)
		return "", false
	}

	return item.Value, true
}

func (s *Store) Delete(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.data[key]
	if exists {
		delete(s.data, key)
		s.expiration.CancelExpiration(key)
		return true
	}
	return false
}

func (s *Store) deleteExpired(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.data, key)
}

func (s *Store) Close() {
	if s.expiration != nil {
		s.expiration.Stop()
	}
}

func (s *Store) SaveToFile(filename string) error {
	return s.persistence.SaveToFile(filename)
}

func (s *Store) LoadFromFile(filename string) error {
	return s.persistence.LoadFromFile(filename)
}
