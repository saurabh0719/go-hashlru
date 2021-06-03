# go-hashlru

A simple thread-safe, fixed size LRU written in Go. Based on [dominictarr's Hashlru Algorithm](https://github.com/dominictarr/hashlru). :arrows_clockwise:

Uses `map[interface{}]interface{}` to allow kv pairs of any type.

Visit `example/example.go` in the root directory for a simple example.

<hr>

## Install

```sh
$ go get github.com/saurabh0719/go-hashlru
```

Latest - `v0.0.3`

<hr>

## API Reference 

HashLRU Struct : 

```go
type HashLRU struct {
	maxSize  int
	size     int
	oldCache map[interface{}]interface{}
	newCache map[interface{}]interface{}
	lock     sync.RWMutex
}
```

### 1. func NewHLRU
```go
func NewHLRU(maxSize int) (*HashLRU, error)
```

Returns a new instance of `HashLRU` of the given size: `maxSize`.

### 2. func (lru *HashLRU) Set
```go
func (lru *HashLRU) Set(key, value interface{})
```

Adds a new key-value pair to the cache.

### 3. func (lru *HashLRU) Get
```go
func (lru *HashLRU) Get(key interface{}) (interface{}, bool)
```

Get the `value` of any `key`. Returns `value, true` if the kv pair is found, else returns `nil, false`.

### 4. func (lru *HashLRU) Has
```go
func (lru *HashLRU) Has(key interface{}) bool
```
Returns `true` if the key exists, else returns `false`.

### 5. func (lru *HashLRU) Remove
```go
func (lru *HashLRU) Remove(key interface{}) bool
```

Deletes the key-value pair and returns `true` if it exists, else returns `false`.

### 6. func (lru *HashLRU) Len
```go
func (lru *HashLRU) Len() int
```

Returns the total number of key-value pairs present in the cache.

### 7. func (lru *HashLRU) Clear
```go
func (lru *HashLRU) Clear()
```

Empties the Cache.

### 8. func (lru *HashLRU) Resize
```go
func (lru *HashLRU) Resize(newSize int) (int, error)
```

Update the `maxSize` of the cache. Items will be evicted automatically to adjust. Returns the *number of evicted key-value pairs* due to the re-size. 