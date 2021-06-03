package lru

import (
	"errors"
	"sync"
)

type HashLRU struct {
	maxSize  int
	size     int
	oldCache map[interface{}]interface{}
	newCache map[interface{}]interface{}
	lock     sync.RWMutex
}

// Returns a new hashlru instance
func NewHLRU(maxSize int) (*HashLRU, error) {

	if maxSize <= 0 {
		return nil, errors.New("Size must be a postive int")
	}

	lru := &HashLRU{
		maxSize:  maxSize,
		size:     0,
		oldCache: make(map[interface{}]interface{}),
		newCache: make(map[interface{}]interface{}),
	}

	return lru, nil

}

func (lru *HashLRU) update(key, value interface{}) {

	lru.newCache[key] = value
	lru.size++

	if lru.size >= lru.maxSize {
		lru.size = 0
		lru.oldCache = lru.newCache
		lru.newCache = make(map[interface{}]interface{})
	}

}

func (lru *HashLRU) Set(key, value interface{}) {

	lru.lock.Lock()

	if _, found := lru.newCache[key]; found {
		lru.newCache[key] = value
	} else {
		lru.update(key, value)
	}

	lru.lock.Unlock()

}

func (lru *HashLRU) Get(key interface{}) (interface{}, bool) {

	lru.lock.Lock()

	if value, found := lru.newCache[key]; found {
		lru.lock.Unlock()
		return value, found
	}

	if value, found := lru.oldCache[key]; found {
		lru.update(key, value)
		lru.lock.Unlock()
		return value, found
	}

	return nil, false

}

// Checks if a key exists in cache
func (lru *HashLRU) Has(key interface{}) bool {

	lru.lock.RLock()

	_, cacheNew := lru.newCache[key]
	_, cacheOld := lru.oldCache[key]

	lru.lock.RUnlock()

	return cacheNew || cacheOld

}

// Removes a key from the cache
func (lru *HashLRU) Remove(key interface{}) (interface{}, bool) {

	lru.lock.Lock()

	if value, found := lru.newCache[key]; found {
		delete(lru.newCache, key)
		lru.lock.Unlock()
		return value, true
	}

	if value, found := lru.oldCache[key]; found {
		delete(lru.oldCache, key)
		lru.lock.Unlock()
		return value, true
	}

	lru.lock.Unlock()

	return nil, false

}

// Returns the number of items in the cache.
func (lru *HashLRU) Len() int {

	lru.lock.RLock()
	length := len(lru.newCache) + len(lru.oldCache)
	lru.lock.RUnlock()

	return length

}

// Clears all entries.
func (lru *HashLRU) Clear() {

	lru.lock.Lock()

	for key, _ := range lru.newCache {
		delete(lru.newCache, key)
	}

	for key, _ := range lru.oldCache {
		delete(lru.oldCache, key)
	}

	lru.lock.Unlock()

}

// Resizes cache, returning number of items deleted
func (lru *HashLRU) Resize(newSize int) (int, error) {

	if newSize <= 0 {
		return 0, errors.New("Size must be a postive int")
	}

	totalItems := len(lru.oldCache) + len(lru.newCache)
	removeCount := totalItems - newSize

	lru.lock.Lock()

	if removeCount < 0 {

		for key, value := range lru.oldCache {
			lru.newCache[key] = value
		}

		lru.oldCache = make(map[interface{}]interface{})
		lru.size = totalItems
		lru.maxSize = newSize
		lru.lock.Unlock()

		return 0, nil

	} else {

		tempKeys := make([]interface{}, 0)
		tempVals := make([]interface{}, 0)

		for key, value := range lru.oldCache {
			tempKeys = append(tempKeys, key)
			tempVals = append(tempVals, value)
		}

		for key, value := range lru.newCache {
			tempKeys = append(tempKeys, key)
			tempVals = append(tempVals, value)
		}

		lru.oldCache = make(map[interface{}]interface{})
		lru.newCache = make(map[interface{}]interface{})
		lru.size = 0
		lru.maxSize = newSize

		for i := 0; i < removeCount; i++ {
			lru.oldCache[tempKeys[i]] = tempVals[i]
		}

		lru.lock.Unlock()

		return removeCount, nil

	}

}
