package lttlru

import (
	"sync"
	"testing"
	"time"
)

func readCache(cache *LruWithTTL, key interface{}) (interface{}, bool) {
	return cache.GetWithTTL(key)
}

func writeCache(cache *LruWithTTL, key, value interface{}, ttl time.Duration) bool {
	return cache.AddWithTTL(key, value, ttl)
}

func TestDataRace(t *testing.T) {
	cache, err := NewTTL(4)
	if err != nil {
		t.Error(err)
	}

	type payload struct {
		Key   interface{}
		Value interface{}
		TTL   time.Duration
	}

	payloads := []payload{
		{"1", 1, time.Minute},
		{"2", 2, time.Minute},
		{"3", 3, time.Minute},
		{"4", 4, time.Minute},
	}

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, p := range payloads {
				writeCache(cache, p.Key, p.Value, p.TTL)
				v, _ := readCache(cache, p.Key)
				t.Log(v)

			}

		}()
	}

	for i := 0; i < 10000; i++ {

		wg.Add(1)
		go func() {
			wg.Done()
			for _, p := range payloads {
				v, _ := readCache(cache, p.Key)
				t.Log(v)
			}
		}()
	}

	//time.Sleep(1 * time.Second)
	wg.Wait()
	for _, p := range payloads {
		v, ok := readCache(cache, p.Key)
		if ok {
			t.Log(v)
		}
	}

}
