// runtime.memhash implementation taken from:
// https://github.com/dgraph-io/ristretto/blob/master/z/rtutil.go
// License can be found:
// https://raw.githubusercontent.com/dgraph-io/ristretto/master/LICENSE

package hashmap

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/cespare/xxhash"
	"github.com/segmentio/fasthash/fnv1a"
)

var (
	// make them var for tests
	loadFactor = 0.75 // if this goes below 0.5 the bounds won't work
	length     = 3
)

type entry struct {
	hash  uint64
	key   string
	value interface{}
}

// Hashmap is a naive implementation of a hashmap struct.
type Hashmap struct {
	lbound int
	ubound int
	length int
	fn     func(string) uint64

	lock    uintptr
	buckets [][8]entry
}

// NewFNV1aHashmap returns a hashmap using the fnv1a
// hashing function.
func NewFNV1aHashmap() *Hashmap {
	return NewHashmap(fnv1a.HashString64)
}

// NewXXHashmap returns a hashmap using the xxhash
// hashing function.
func NewXXHashmap() *Hashmap {
	return NewHashmap(xxhash.Sum64String)
}

// NewHashmap creates a new, empty, hashmap.
func NewHashmap(fn func(string) uint64) *Hashmap {
	return newWithCap(fn, 1<<length)
}

func newWithCap(fn func(string) uint64, cap int) *Hashmap {
	h := &Hashmap{
		lbound:  int(float64(int(1)<<length) * (1 - loadFactor)),
		ubound:  int(float64(int(1)<<length) * loadFactor),
		length:  cap,
		buckets: make([][8]entry, cap),
		fn:      fn,
	}

	if h.lbound == 1<<length || h.ubound == 1<<length {
		panic("invalid load factor")
	}

	return h
}

// Add inserts the value v associated with the key k into the hashmap.
// Redistribution of keys occurs if load factor is surpassed.
func (h *Hashmap) Add(k string, v interface{}) bool {
	for {
		if atomic.CompareAndSwapUintptr(&h.lock, 0, 1) {
			break
		}
	}

	hash := h.fn(k)
	idx := hash & (uint64(h.length) - 1)

	var (
		exists bool
		length int
		target = -1
	)
	for i := range h.buckets[idx] {
		if h.buckets[idx][i].hash == hash {
			h.buckets[idx][i].value = v
			exists = true
		}

		if h.buckets[idx][i].hash <= 0 {
			if target == -1 {
				target = i
			}
		} else {
			length++
		}
	}

	if !exists {
		if target == -1 {
			target = length
		}
		h.buckets[idx][target].hash = hash
		h.buckets[idx][target].key = k
		h.buckets[idx][target].value = v
	}

	// Assume even distribution
	if length >= h.ubound {
		h.length *= 2
		h.resize()
	}

	atomic.StoreUintptr(&h.lock, 0)
	return !exists
}

// Delete removes the key from the map, if it exists,
// and returns whether or not it was deleted.
func (h *Hashmap) Delete(k string) bool {
	for {
		if atomic.CompareAndSwapUintptr(&h.lock, 0, 1) {
			break
		}
	}

	hash := h.fn(k)
	idx := hash & (uint64(h.length) - 1)

	var (
		exists bool
		length int
	)
	for i := range h.buckets[idx] {
		if h.buckets[idx][i].hash == hash {
			h.buckets[idx][i].hash = 0
			h.buckets[idx][i].key = ""
			h.buckets[idx][i].value = nil
			exists = true
		}

		if h.buckets[idx][i].hash > 0 {
			length++
		}
	}

	if length <= h.lbound {
		h.length /= 2
		h.resize()
	}

	atomic.StoreUintptr(&h.lock, 0)
	return exists
}

func (h *Hashmap) resize() {
	buckets := make([][8]entry, h.length)

	for i := range h.buckets {
		var length int
		for j := range h.buckets[i] {
			if h.buckets[i][j].hash == 0 {
				continue
			}
			idx := h.buckets[i][j].hash & (uint64(h.length) - 1)

			buckets[idx][length] = h.buckets[i][j]
			length++
		}
	}

	h.buckets = buckets
}

// Lookup will try to retrieve the value associated with
// the specified key.
func (h *Hashmap) Lookup(k string) (interface{}, bool) {
	for {
		if atomic.CompareAndSwapUintptr(&h.lock, 0, 1) {
			break
		}
	}

	hash := h.fn(k)
	idx := hash & (uint64(h.length) - 1)

	for i := range h.buckets[idx] {
		if h.buckets[idx][i].hash == hash {
			atomic.StoreUintptr(&h.lock, 0)
			return h.buckets[idx][i].value, true
		}
	}

	atomic.StoreUintptr(&h.lock, 0)
	return nil, false
}

func spew(buckets [][8]entry) string {
	var lengths []int
	for i := range buckets {
		var length int
		for j := range buckets[i] {
			if buckets[i][j].hash != 0 {
				length++
			}
		}
		lengths = append(lengths, int(length))
	}

	var b strings.Builder
	for i := range lengths {
		b.WriteString(fmt.Sprintf("[%d]", lengths[i]))
	}
	return b.String()
}
