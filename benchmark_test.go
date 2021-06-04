package hlru

import (
	"math/rand"
	"testing"
)

func BenchmarkHLRU_Rand(b *testing.B) {

	lru, _ := NewHLRU(8192)

	trace := make([]int64, b.N*2)

	for i := 0; i < b.N*2; i++ {
		trace[i] = rand.Int63() % 32768
	}

	b.ResetTimer()

	var hit, miss int

	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			lru.Set(trace[i], trace[i])
		} else {
			_, ok := lru.Get(trace[i])
			if ok {
				hit++
			} else {
				miss++
			}
		}
	}

	b.Logf("Hit: %d Miss: %d Ratio: %f", hit, miss, float64(hit)/float64(miss))

}

func BenchmarkHLRU_Freq(b *testing.B) {

	lru, _ := NewHLRU(8192)

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = rand.Int63() % 16384
		} else {
			trace[i] = rand.Int63() % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lru.Set(trace[i], trace[i])
	}

	var hit, miss int

	for i := 0; i < b.N; i++ {
		_, ok := lru.Get(trace[i])
		if ok {
			hit++
		} else {
			miss++
		}
	}

	b.Logf("Hit: %d Miss: %d Ratio: %f", hit, miss, float64(hit)/float64(miss))

}
