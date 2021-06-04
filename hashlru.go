package hlru

import (
	"errors"
	"sync"
	"math"
)

/*
The HashLRU algorithm maintains two separate maps 
and bulk eviction happens only after both the maps fill up

Hence the the callBack function is triggered in bulk and 
is not an accurate measure. Use NewWithEvict() with caution.
*/

type HashLRU struct {
	maxSize  					int
	size     					int
	oldCache, newCache 			map[interface{}]interface{}
	onEvictedCB					func (key, value interface{})
	lock     					sync.RWMutex
}

type KVPair struct {
	key, value			interface{}
}

// Returns a new hashlru instance
func NewHLRU(maxSize int) (*HashLRU, error) {

	return NewWithEvict(maxSize, nil)

}

func NewWithEvict(maxSize int, onEvict func(key, value interface{})) (*HashLRU, error) {

	if maxSize <= 0 {
		return nil, errors.New("Size must be a postive int")
	}

	lru := &HashLRU{
		maxSize:  maxSize,
		size:     0,
		onEvictedCB: onEvict,
		oldCache: make(map[interface{}]interface{}),
		newCache: make(map[interface{}]interface{}),
	}

	return lru, nil

}

/*
update(key, value interface{}) is used internally in Get() and Set()
to impose least recently by pushing all recently accessed keys to the newCache
and the oldCache acts as a back up dump once newCache fills up. 

Bulk eviction takes place from the oldCache
*/

func (lru *HashLRU) update(key, value interface{}) {

	lru.newCache[key] = value
	lru.size++

	if lru.size >= lru.maxSize {
		lru.size = 0

		if lru.onEvictedCB != nil {
			for key, value := range lru.oldCache {
				lru.onEvictedCB(key, value)
			}
		}
		
		lru.oldCache = make(map[interface{}]interface{})
		for key, value := range lru.newCache {
			lru.oldCache[key] = value
		}

		lru.newCache = make(map[interface{}]interface{})
	}

}

// Set a value and update the cache
func (lru *HashLRU) Set(key, value interface{}) {

	lru.lock.Lock()

	if _, found := lru.newCache[key]; found {
		lru.newCache[key] = value
	} else {
		lru.update(key, value)
	}

	lru.lock.Unlock()

}

// Get a value and update the cache
func (lru *HashLRU) Get(key interface{}) (interface{}, bool) {

	lru.lock.Lock()

	if value, found := lru.newCache[key]; found {
		lru.lock.Unlock()
		return value, found
	}

	if value, found := lru.oldCache[key]; found {
		delete(lru.oldCache, key)
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

	if val, found := lru.newCache[key]; found {
		delete(lru.newCache, key)
		lru.size--
		if lru.onEvictedCB != nil {
			lru.onEvictedCB(key, val)
		}
		lru.lock.Unlock()
		return true
	}

	if val, found := lru.oldCache[key]; found {
		delete(lru.oldCache, key)
		if lru.onEvictedCB != nil {
			lru.onEvictedCB(key, val)
		}
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

	if lru.onEvictedCB != nil {
		for key, value := range lru.oldCache {
			lru.onEvictedCB(key, value)
		}
		for key, value := range lru.newCache {
			lru.onEvictedCB(key, value)
		}
	}

	lru.oldCache = make(map[interface{}]interface{})
	lru.newCache = make(map[interface{}]interface{})
	lru.size = 0

	lru.lock.Unlock()

}

func (lru* HashLRU) Keys() []interface{} {

	lru.lock.RLock()

	tempKeys := make([]interface{}, 0)

	for key, _ := range lru.oldCache {
		// tempKeys = append(tempKeys, key)
		if _, found := lru.newCache[key]; !found {
			tempKeys = append(tempKeys, key)
		}
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

	for key, value := range lru.oldCache {
		if _, found := lru.newCache[key]; !found {
			tempVals = append(tempVals, value)
		}
	}

	for _, value := range lru.newCache {
		tempVals = append(tempVals, value)
	}
	
	lru.lock.RUnlock()
	return tempVals

}

func (lru *HashLRU) all() []*KVPair {

	lru.lock.RLock()

	allPairs := []*KVPair{}

	for key, value := range lru.oldCache {
		if _, found := lru.newCache[key]; !found {

			kvPair := new(KVPair)
			kvPair.key = key
			kvPair.value = value
			allPairs = append(allPairs, kvPair)
		}
	}

	for key, value := range lru.newCache {
		kvPair := new(KVPair)
		kvPair.key = key
		kvPair.value = value
		allPairs = append(allPairs, kvPair)
	}
	
	lru.lock.RUnlock()
	return allPairs

}

// Resizes cache, returning number of items deleted
func (lru *HashLRU) Resize(newSize int) (int, error) {

	if newSize <= 0 {
		return 0, errors.New("Size must be a postive int")
	}

	totalItems := lru.Len()
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

		allPairs := lru.all()

		lru.lock.Lock()

		lru.oldCache = make(map[interface{}]interface{})
		lru.newCache = make(map[interface{}]interface{})
		lru.size = 0
		lru.maxSize = newSize

		var i = 0

		for i < removeCount {
			if lru.onEvictedCB != nil {
				lru.onEvictedCB(allPairs[i].key, allPairs[i].value)
			}
			i++
		}

		for i < len(allPairs) {
			lru.oldCache[allPairs[i].key] = allPairs[i].value
			i++
		}

		lru.lock.Unlock()

		return removeCount, nil

	}

}
