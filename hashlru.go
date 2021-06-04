package hlru

import (
	"errors"
	"sync"
	"math"
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

	lru.lock.Unlock()
	return nil, false

}

// Peek the value of a key without updating the cache
func (lru *HashLRU) Peek(key interface{}) (interface{}, bool) {

	lru.lock.RLock()

	if value, found := lru.newCache[key]; found {
		lru.lock.RUnlock()
		return value, found
	}

	if value, found := lru.oldCache[key]; found {
		lru.lock.RUnlock()
		return value, found
	}

	lru.lock.RUnlock()
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
func (lru *HashLRU) Remove(key interface{}) bool {

	lru.lock.Lock()

	if _, found := lru.newCache[key]; found {
		delete(lru.newCache, key)
		lru.lock.Unlock()
		return true
	}

	if _, found := lru.oldCache[key]; found {
		delete(lru.oldCache, key)
		lru.lock.Unlock()
		return true
	}

	lru.lock.Unlock()

	return false

}

// Returns the number of items in the cache.
func (lru *HashLRU) Len() int {

	lru.lock.RLock()
	
	if lru.size == 0 {
		lru.lock.RUnlock()
		return len(lru.oldCache)
	}

	oldCacheSize := 0

	for key, _ := range lru.oldCache {
		if _, found := lru.newCache[key]; !found {
			oldCacheSize++
		}
	}

	lru.lock.RUnlock()
	return int(math.Min(float64(lru.size + oldCacheSize), float64(lru.maxSize)))

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

func (lru* HashLRU) Keys() []interface{} {

	lru.lock.RLock()

	tempKeys := make([]interface{}, 0)

	for key, _ := range lru.oldCache {
		tempKeys = append(tempKeys, key)
	}

	for key, _ := range lru.newCache {
		tempKeys = append(tempKeys, key)
	}
	
	lru.lock.RUnlock()
	return tempKeys

}

func (lru *HashLRU) Vals() []interface{} {

	lru.lock.RLock()

	tempVals := make([]interface{}, 0)

	for _, value := range lru.oldCache {
		tempVals = append(tempVals, value)
	}

	for _, value := range lru.newCache {
		tempVals = append(tempVals, value)
	}
	
	lru.lock.RUnlock()
	return tempVals

}

// Resizes cache, returning number of items deleted
func (lru *HashLRU) Resize(newSize int) (int, error) {

	if newSize <= 0 {
		return 0, errors.New("Size must be a postive int")
	}

	totalItems := len(lru.oldCache) + len(lru.newCache)
	removeCount := totalItems - newSize

	if removeCount < 0 {

		lru.lock.Lock()

		for key, value := range lru.oldCache {
			lru.newCache[key] = value
		}

		lru.oldCache = make(map[interface{}]interface{})
		lru.size = totalItems
		lru.maxSize = newSize
		lru.lock.Unlock()

		return 0, nil

	} else {

		tempKeys := lru.Keys()
		tempVals := lru.Vals()

		lru.lock.Lock()

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
