package cache_test

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/Karanth1r3/l_0/internal/cache"
	"github.com/Karanth1r3/l_0/internal/model"
)

func TestCache(t *testing.T) {
	cache := cache.NewCache()
	data := []byte("trash")
	orderUID := "qwerty"
	tests := []struct {
		name     string
		isOk     bool
		orderUID string
		err      error
	}{
		{
			name:     "ok",
			isOk:     true,
			orderUID: orderUID,
		},
		{
			name:     "not_found",
			isOk:     false,
			orderUID: "trash",
			err:      model.ErrNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache.Write(orderUID, data)
			received, err := cache.Read(test.orderUID)
			if err != nil {
				if test.isOk {
					t.Fatal("unexpected behaviour")
				}
				return
			}
			if !bytes.Equal(received, data) {
				t.Fatal("unexpected behaviour")
			}
		})
	}
}

func TestCacheParallel(t *testing.T) {
	cache := cache.NewCache()
	wg := &sync.WaitGroup{}
	go read(cache, wg)
	go write(cache, wg)
	wg.Wait()
}

func read(c *cache.Cache, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for i := 0; i < 100; i++ {
		_, _ = c.Read(fmt.Sprintf("uid %d", i))
	}
}

func write(c *cache.Cache, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for i := 0; i < 100; i++ {
		c.Write(fmt.Sprintf("uid %d", i), []byte("trash2"))
	}
}
