# go-hashlru

![](https://img.shields.io/github/go-mod/go-version/saurabh0719/go-hashlru) ![](https://img.shields.io/github/v/release/saurabh0719/go-hashlru?color=FFD500)

A simple thread-safe, fixed size LRU written in Go. Based on [dominictarr's Hashlru Algorithm](https://github.com/dominictarr/hashlru). :arrows_clockwise:

Uses `map[interface{}]interface{}` to allow kv pairs of any type. The `hlru` package contains all the necessary functions.

<div align="center">
    <strong><a href="https://github.com/saurabh0719/go-hashlru">Github</a> | <a href="https://saurabh0719.github.io">Website</a> | <a href="https://github.com/saurabh0719/go-hashlru/releases">Releases</a> </strong>
</div>

```go
cache, _ := hlru.NewHLRU(100)

cache.Set("key", "value")

val, _ := cache.Get("key")

fmt.Println(val)
// value
```

Visit `example/example.go` in the root directory for a simple example.

<hr>

<span id="contents"></span>

### Table of Contents :
* [Installation](#installation)
* [API reference](#api)
    * [HashLRU struct](#type)
    * [Functions](#func)
* [Tests](#tests)
* [Benchmark](#bench)

<span id="installation"></span>
### Installation

```sh
$ go get github.com/saurabh0719/go-hashlru
```

Latest - `v0.1.0`

<hr>

### API Reference 

<span id="type"></span>
#### HashLRU Struct  

```go
type HashLRU struct {
	maxSize            int
	size               int
	oldCache, newCache map[interface{}]interface{}
	onEvictedCB        func(key, value interface{})
	lock               sync.RWMutex
}
```

The HashLRU algorithm maintains two separate maps
and bulk eviction happens only after both the maps fill up. Hence `onEvictedCB` is triggered in bulk and is not an accurate measure of timely LRU-ness.

As explained by [dominictarr](https://github.com/dominictarr/hashlru) :
* This algorithm does not give you an ordered list of the N most recently used items, but has the important properties of the LRU (bounded memory use and O(1) time complexity)


This implementation uses `sync.RWMutex` to ensure thread-safety.

[Go back to the table of contents](#contents)

<span id="func"></span>

#### func NewHLRU
```go
func NewHLRU(maxSize int) (*HashLRU, error)
```

Returns a new instance of `HashLRU` of the given size: `maxSize`.

#### func NewWithEvict
```go
func NewWithEvict(maxSize int, onEvict func(key, value interface{})) (*HashLRU, error)
```

Takes `maxSzie` and a callback function `onEvict` as arguments and returns a new instance of `HashLRU` of the given size.

#### func (lru *HashLRU) Set
```go
func (lru *HashLRU) Set(key, value interface{})
```

Adds a new key-value pair to the cache *and* updates it.

#### func (lru *HashLRU) Get
```go
func (lru *HashLRU) Get(key interface{}) (interface{}, bool)
```

Get the `value` of any `key` *and* updates the cache. Returns `value, true` if the kv pair is found, else returns `nil, false`.

#### func (lru *HashLRU) Has
```go
func (lru *HashLRU) Has(key interface{}) bool
```
Returns `true` if the key exists, else returns `false`. 

#### func (lru *HashLRU) Remove
```go
func (lru *HashLRU) Remove(key interface{}) bool
```

Deletes the key-value pair and returns `true` if its successful, else returns `false`.

[Go back to the table of contents](#contents)

#### func (lru *HashLRU) Peek
```go
func (lru *LRU) Peek(key interface{}) (interface{}, bool)
```

Get the value of a key without updating the cache. Returns `value, true` if the kv pair is found, else returns `nil, false`. 

#### func (lru *HashLRU) Clear
```go
func (lru *HashLRU) Clear()
```

Empties the Cache.

#### func (lru *HashLRU) Resize
```go
func (lru *HashLRU) Resize(newSize int) (int, error)
```

Update the `maxSize` of the cache. Items will be evicted automatically to adjust. Returns the *number of evicted key-value pairs* due to the re-size. 

#### func (lru *HashLRU) Len
```go
func (lru *HashLRU) Len() int
```

Returns the total number of key-value pairs present in the cache.

#### func (lru *HashLRU) Keys
```go
func (lru* HashLRU) Keys() []interface{}
```

Returns a slice of all the Keys in the cache.

#### func (lru *HashLRU) Vals
```go
func (lru *HashLRU) Vals() []interface{}
```

Returns a slice of all the Values in the cache.

[Go back to the table of contents](#contents)

<hr>

<span id="test"></span>

### Tests
```sh
$ go test 
```

Use `-v` for an detailed output.

<hr>

<span id="func"></span>

### Benchmark 
```sh
$ go test -run=XXX -bench=.
```

HashLRU has a better hit/miss ratio for `BenchmarkHLRU_Rand` and `BenchmarkHLRU_Freq` (negligible) as compared to same benchmark tests from [golang-lru](https://github.com/hashicorp/golang-lru). However, `golang-lru` does have a slightly lower `ns/op`.

[Go back to the table of contents](#contents)

<hr>
