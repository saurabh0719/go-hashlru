package hlru

import (
	"testing"
	// "fmt"
)

func Test_HLRU(t *testing.T) {
	
	lru, err := NewHLRU(100)

	if err != nil {
		t.Fatalf("Error in creating LRU: %v", err)
	}

	for i := 0; i < 150; i++ {
		lru.Set(i, i)
	}

	if lru.Len() != 100 {
		t.Fatalf("Error in LRU length: %v", lru.Len())
	}

	lru.Clear() 

	if lru.Len() != 0 {
		t.Fatalf("Error in LRU Clear(): %v", lru.Len())
	}

	lru.Set(1, 1)
	lru.Set(2, 1)

	keys := lru.Keys()

	for i:=0; i< len(keys); i++ {
		_, ok := lru.Get(keys[i]) 
		if !ok {
			t.Fatalf("Error: %v", keys[i])
		}
	}

	if len(keys) != lru.Len() {
		t.Fatalf("Error: %v", keys)
	}

}

func Test_onEvict(t *testing.T) {

	evicted := 0
	onEvict := func(key, value interface{}) {
		if key != value {
			t.Fatalf("Evict values not equal (%v!=%v)", key, value)
		}
		evicted++
	}

	lru, err := NewWithEvict(100, onEvict)

	if err != nil {
		t.Fatalf("Error in creating LRU: %v", err)
	}

	for i := 0; i < 150; i++ {
		lru.Set(i, i)
	}

	if lru.Len() != 100 {
		t.Fatalf("Error in LRU length: %v", lru.Len())
	}

	keys := lru.Keys()
	// vals := lru.Vals()

	for i:=0; i<len(keys); i++ {
		if lru.Has(keys[i]) != true {
			t.Fatalf("Error in Has() Keys()")
		}
	}

	for i:=0; i< lru.Len(); i++ {
		_, ok := lru.Peek(i) 
		if !ok {
			t.Fatalf("Error in Peek()")
		}
	}

	if evicted != 0 {
		t.Fatalf("Error in evict callback: %v", evicted)
	}

}

func Test_Remove_Resize(t *testing.T) {

	lru, _ := NewHLRU(2)

	lru.Set(1, 1)
	lru.Set(2, 2)

	ok := lru.Remove(2)

	if !ok {
		t.Fatalf("Error in Remove()")
	}

	if lru.Has(2) != false {
		t.Fatalf("Error in Remove()")
	}

	lru.Set(3, 3)
	lru.Set(2, 2)

	if lru.Has(2) == false {
		t.Fatalf("Error in Has()")
	}

	if lru.Has(1) != false {
		t.Fatalf("Error in Remove()")
	}

	lru.Clear()

	lru.Set(1, 1)
	lru.Set(2, 2)

	var evicted, _ = lru.Resize(1)

	if evicted != 1 {
		t.Fatalf("Error in Down Sizing")
	}

	// var keys = lru.Keys()
	// // vals := lru.Vals()

	// for i:=0; i<len(keys); i++ {
	// 	fmt.Println("Key first", keys[i])
	// }

	evicted, _ = lru.Resize(2)

	if evicted != 0 {
		t.Fatalf("Error in Down Sizing")
	}

	lru.Set(3,3)
	lru.Set(4,4)
	lru.Set(5,5)

	if lru.Len() != 2 {
		t.Fatalf("Error in LRU length: %v", lru.Len())
	}

	// keys = lru.Keys()
	// // vals := lru.Vals()

	// for i:=0; i<len(keys); i++ {
	// 	fmt.Println("Key", keys[i])
	// }

	// fmt.Print("LEN ", lru.Len())
}