package store

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type PersistenceManager struct {
	store *Store
}

func NewPersistenceManager(store *Store) *PersistenceManager {
	return &PersistenceManager{
		store: store,
	}
}

func (pm *PersistenceManager) SaveToFile(filename string) error {
	pm.store.mutex.RLock()
	defer pm.store.mutex.RUnlock()

	data := make(map[string]*Item)
	for k, v := range pm.store.data {
		if v.HasTTL && time.Now().After(v.ExpiresAt) {
			continue
		}
		data[k] = v
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (pm *PersistenceManager) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var data map[string]*Item
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("failed to decode file: %v", err)
	}

	pm.store.mutex.Lock()
	defer pm.store.mutex.Unlock()

	pm.store.data = make(map[string]*Item)

	now := time.Now()
	for key, item := range data {
		if item.HasTTL && now.After(item.ExpiresAt) {
			continue
		}

		pm.store.data[key] = item

		if item.HasTTL {
			remaining := item.ExpiresAt.Sub(now)
			if remaining > 0 {
				pm.store.expiration.SetExpiration(key, remaining)
			}
		}
	}

	return nil
}

func (pm *PersistenceManager) SaveSnapshot() error {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("kvstore_snapshot_%s.json", timestamp)
	return pm.SaveToFile(filename)
}

func (pm *PersistenceManager) GetStats() map[string]interface{} {
	pm.store.mutex.RLock()
	defer pm.store.mutex.RUnlock()

	totalKeys := len(pm.store.data)
	expiredKeys := 0
	now := time.Now()

	for _, item := range pm.store.data {
		if item.HasTTL && now.After(item.ExpiresAt) {
			expiredKeys++
		}
	}

	return map[string]interface{}{
		"total_keys":   totalKeys,
		"expired_keys": expiredKeys,
		"active_keys":  totalKeys - expiredKeys,
	}
}
