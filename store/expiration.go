package store

import (
	"sync"
	"time"
)

type ExpirationManager struct {
	store   *Store
	timers  map[string]*time.Timer
	mutex   sync.RWMutex
	stopped bool
}

func NewExpirationManager(store *Store) *ExpirationManager {
	return &ExpirationManager{
		store:  store,
		timers: make(map[string]*time.Timer),
	}
}

func (em *ExpirationManager) SetExpiration(key string, ttl time.Duration) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if em.stopped {
		return
	}

	if timer, exists := em.timers[key]; exists {
		timer.Stop()
	}

	timer := time.AfterFunc(ttl, func() {
		em.store.deleteExpired(key)
		em.mutex.Lock()
		delete(em.timers, key)
		em.mutex.Unlock()
	})

	em.timers[key] = timer
}

func (em *ExpirationManager) CancelExpiration(key string) {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	if timer, exists := em.timers[key]; exists {
		timer.Stop()
		delete(em.timers, key)
	}
}

func (em *ExpirationManager) Stop() {
	em.mutex.Lock()
	defer em.mutex.Unlock()

	em.stopped = true
	for _, timer := range em.timers {
		timer.Stop()
	}
	em.timers = make(map[string]*time.Timer)
}
