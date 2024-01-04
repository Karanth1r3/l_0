package cache

import (
	"sync"

	"github.com/Karanth1r3/l_0/internal/model"
)

type Cache struct {
	records   map[string][]byte
	recordsMu sync.RWMutex
}

// Constructor
func NewCache() *Cache {
	return &Cache{
		records: make(map[string][]byte),
	}
}

// Write record to cache
func (c *Cache) Write(orderUID string, data []byte) {
	c.recordsMu.Lock() // to avoid errors while writing to map
	defer c.recordsMu.Unlock()
	c.records[orderUID] = data
}

// Read record from cache by orderUID
func (c *Cache) Read(orderUID string) ([]byte, error) {
	c.recordsMu.RLock() // to avoid errors linked to concurrent w/r to map (map is not thread safe)
	defer c.recordsMu.RUnlock()
	r, ok := c.records[orderUID]
	if !ok {
		return nil, model.ErrNotFound
	}

	return r, nil
}
