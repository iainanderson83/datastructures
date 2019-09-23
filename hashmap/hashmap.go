// runtime.memhash implementation taken from:
// https://github.com/dgraph-io/ristretto/blob/master/z/rtutil.go
// License can be found:
// https://raw.githubusercontent.com/dgraph-io/ristretto/master/LICENSE

package hashmap

import (
	"math"
	"strings"
	"unsafe"

	"github.com/segmentio/fasthash/fnv1a"
)

type stringStruct struct {
	str unsafe.Pointer
	len int
}

var (
	bucketSize float64 = 3
)

type tuple struct {
	key   string
	value interface{}
}

// Hashmap is a naive implementation of a hashmap struct.
type Hashmap struct {
	loadFactor float64
	buckets    [][]tuple
	fn         func(string) uint64
}

// NewFNV1aHashmap returns a hashmap using the fnv1a
// hashing function.
func NewFNV1aHashmap() *Hashmap {
	return NewHashmap(fnv1a.HashString64)
}

// NewRuntimeHashmap returns a hashmap using the runtime.memhash
// hashing function.
func NewRuntimeHashmap() *Hashmap {
	return NewHashmap(memHashString)
}

// NewHashmap creates a new, empty, hashmap.
func NewHashmap(fn func(string) uint64) *Hashmap {
	return &Hashmap{
		loadFactor: 0.75,
		buckets:    make([][]tuple, int(math.Pow(2, bucketSize))),
		fn:         fn,
	}
}

// Add inserts the value v associated with the key k into the hashmap.
// Redistribution of keys occurs if load factor is surpassed.
func (h *Hashmap) Add(k string, v interface{}) {
	k = strings.ToLower(k)
	hash := h.fn(k)
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

		for i := range h.buckets {
			for j := range h.buckets[i] {
				key := strings.ToLower(h.buckets[i][j].key)
				hash := h.fn(key)
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
	}
}

func isPowof2(n int) bool {
	return n != 0 && (n&(n-1)) == 0
}

// Lookup will try to retrieve the value associated with
// the specified key.
func (h *Hashmap) Lookup(k string) interface{} {
	k = strings.ToLower(k)
	hash := h.fn(k)
	idx := hash & uint64(len(h.buckets)-1)

	for _, tpl := range h.buckets[idx] {
		if tpl.key == k {
			return tpl.value
		}
	}

	return nil
}

//go:noescape
//go:linkname memhash runtime.memhash
func memhash(p unsafe.Pointer, h, s uintptr) uintptr

// MemHash is the hash function used by go map, it utilizes available hardware instructions(behaves
// as aeshash if aes instruction is available).
// NOTE: The hash seed changes for every process. So, this cannot be used as a persistent hash.
func memHash(data []byte) uint64 {
	ss := (*stringStruct)(unsafe.Pointer(&data))
	return uint64(memhash(ss.str, 0, uintptr(ss.len)))
}

// MemHashString is the hash function used by go map, it utilizes available hardware instructions
// (behaves as aeshash if aes instruction is available).
// NOTE: The hash seed changes for every process. So, this cannot be used as a persistent hash.
func memHashString(str string) uint64 {
	ss := (*stringStruct)(unsafe.Pointer(&str))
	return uint64(memhash(ss.str, 0, uintptr(ss.len)))
}
