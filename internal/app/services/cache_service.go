package services

import (
	"sync"
)

// CacheService is a simple in-memory key-value store.
// It is safe for concurrent access using read-write mutexes.
type CacheService struct {
	mu   sync.RWMutex           // Mutex to ensure thread-safe access (read/write)
	data map[string]interface{} // Underlying map to store cached data
}

// InitCacheService initializes and returns a new CacheService instance.
// It creates an empty map ready to store data.
func InitCacheService() *CacheService {
	return &CacheService{
		data: make(map[string]interface{}),
	}
}

// Get retrieves the value associated with the given key.
// Returns the value and a boolean indicating if the key exists.
// Uses a read lock to allow multiple concurrent readers.
func (service *CacheService) Get(key string) (interface{}, bool) {
	service.mu.RLock()         // Acquire read lock
	defer service.mu.RUnlock() // Release read lock when function exits

	val, ok := service.data[key] // Look up key in the map
	return val, ok
}

// Set stores or updates the value for the given key.
// Uses a write lock to ensure only one goroutine can write at a time.
func (service *CacheService) Set(key string, value interface{}) {
	service.mu.Lock()         // Acquire write lock
	defer service.mu.Unlock() // Release write lock when function exits

	service.data[key] = value // Set or overwrite the key-value pair
}
