package hashmap

import (
	"fmt"
	"math"
	"strings"

	"github.com/segmentio/fasthash/fnv1a"
)

var (
	bucketSize float64 = 3
	debug      bool
)

type tuple struct {
	key   string
	value interface{}
}

// Hashmap is a naive implementation of a hashmap struct.
type Hashmap struct {
	loadFactor float64
	buckets    [][]tuple
}

// NewHashmap creates a new, empty, hashmap.
func NewHashmap() *Hashmap {
	return &Hashmap{
		loadFactor: 0.75,
		buckets:    make([][]tuple, int(math.Pow(2, bucketSize))),
	}
}

// Add inserts the value v associated with the key k into the hashmap.
// Redistribution of keys occurs if load factor is surpassed.
func (h *Hashmap) Add(k string, v interface{}) {
	k = strings.ToLower(k)
	hash := fnv1a.HashString64(k)
	idx := hash & uint64(len(h.buckets)-1)

	var added bool
	for i, tpl := range h.buckets[idx] {
		if tpl.key == k {
			h.buckets[idx][i] = tuple{k, v}
			added = true
			break
		}
	}

	if !added {
		h.buckets[idx] = append(h.buckets[idx], tuple{k, v})
	}

	// Assume even distribution
	if float64(len(h.buckets[idx]))/math.Pow(2, bucketSize) > h.loadFactor {
		// Redistribute
		newBuckets := make([][]tuple, int(math.Pow(2, math.Log2(float64(len(h.buckets)))+1)))

		if !isPowof2(len(newBuckets)) {
			panic("invalid number of buckets")
		}

		print("before: %+v", h.buckets)

		for i := range h.buckets {
			for j := range h.buckets[i] {
				key := strings.ToLower(h.buckets[i][j].key)
				hash := fnv1a.HashString64(key)
				idx := hash & uint64(len(newBuckets)-1)

				var added bool
				for l, tpl := range newBuckets[idx] {
					if tpl.key == key {
						newBuckets[idx][l] = tuple{key, h.buckets[i][j].value}
						added = true
						break
					}
				}

				if !added {
					newBuckets[idx] = append(newBuckets[idx], tuple{key, h.buckets[i][j].value})
				}
			}
		}
		h.buckets = newBuckets

		print("after: %+v", h.buckets)
	}
}

func print(msg string, args ...interface{}) {
	if debug {
		fmt.Printf(msg+"\n", args...)
	}
}

func isPowof2(n int) bool {
	return n != 0 && (n&(n-1)) == 0
}

// Lookup will try to retrieve the value associated with
// the specified key.
func (h *Hashmap) Lookup(k string) interface{} {
	k = strings.ToLower(k)
	hash := fnv1a.HashString64(k)
	idx := hash & uint64(len(h.buckets)-1)

	for _, tpl := range h.buckets[idx] {
		if tpl.key == k {
			return tpl.value
		}
	}

	return nil
}
