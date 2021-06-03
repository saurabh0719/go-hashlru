package main

import (
	"fmt"
	"github.com/saurabh0719/go-hashlru"
)

func main() {

	cache, _ := lru.NewHLRU(120)

	cache.Set("key", 30)
	cache.Set(20, 5)

	val, _ := cache.Get("key")
	fmt.Println(val)
	// 30

	numVal, _ := cache.Get(20)
	fmt.Println(numVal)
	// 5

	fmt.Println(cache.Has(20))
	// true

}
