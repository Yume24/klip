package hls

import (
	"bytes"
	"sync"
)

type keyCache struct {
	keys map[string][]byte
	lock sync.RWMutex
}

func createKeyCache() *keyCache {
	return &keyCache{keys: make(map[string][]byte)}
}

func (cache *keyCache) getOrFetch(keyURI string) ([]byte, error) {
	cache.lock.RLock()
	cacheValue, ok := cache.keys[keyURI]
	cache.lock.RUnlock()

	if ok {
		return cacheValue, nil
	}

	cache.lock.Lock()
	defer cache.lock.Unlock()

	cacheValue, ok = cache.keys[keyURI]
	if ok {
		return cacheValue, nil
	}

	fetchedKey := &bytes.Buffer{}
	err := getResponseBody(keyURI, fetchedKey)
	if err != nil {
		return nil, err
	}

	cache.keys[keyURI] = fetchedKey.Bytes()

	return fetchedKey.Bytes(), nil
}
