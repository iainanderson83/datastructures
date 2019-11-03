package hashmap

import "sync/atomic"

// OrderedMap is an ordered variant of Hashmap.
type OrderedMap struct {
	lock uintptr
	i    []string
	m    *Hashmap
}

// NewOrderedMap creates a new ordered map with the specified hashing function.
func NewOrderedMap(fn func(string) uint64) *OrderedMap {
	return &OrderedMap{m: NewHashmap(fn)}
}

// Iter calls the specified cb for each key/value pair in the map
// in the inserted order.
func (o *OrderedMap) Iter(fn func(k string, v interface{}) bool) {
	for {
		if atomic.CompareAndSwapUintptr(&o.lock, 0, 1) {
			break
		}
	}

	for i := range o.i {
		v, b := o.m.Lookup(o.i[i])
		if !b {
			o.i = append(o.i[:i], o.i[i+1:]...)
			continue
		}

		if !fn(o.i[i], v) {
			atomic.StoreUintptr(&o.lock, 0)
			return
		}
	}

	atomic.StoreUintptr(&o.lock, 0)
}

// Lookup returns the value associated with the specified key in the map.
func (o *OrderedMap) Lookup(k string) (interface{}, bool) {
	for {
		if atomic.CompareAndSwapUintptr(&o.lock, 0, 1) {
			break
		}
	}

	v, b := o.m.Lookup(k)

	atomic.StoreUintptr(&o.lock, 0)
	return v, b
}

// Delete removes the value associated with the specified key from the map.
func (o *OrderedMap) Delete(k string) bool {
	for {
		if atomic.CompareAndSwapUintptr(&o.lock, 0, 1) {
			break
		}
	}

	deleted := o.m.Delete(k)
	if deleted {
		idx := -1
		for i := range o.i {
			if o.i[i] == k {
				idx = i
				break
			}
		}

		if idx >= 0 {
			o.i = append(o.i[:idx], o.i[idx+1:]...)
		}
	}

	atomic.StoreUintptr(&o.lock, 0)
	return deleted
}

// Add adds the specified value to the map with the specified key.
func (o *OrderedMap) Add(k string, v interface{}) bool {
	for {
		if atomic.CompareAndSwapUintptr(&o.lock, 0, 1) {
			break
		}
	}

	added := o.m.Add(k, v)
	if added {
		o.i = append(o.i, k)
	}

	atomic.StoreUintptr(&o.lock, 0)
	return added
}

// Len returns the number of elements in  the map.
func (o *OrderedMap) Len() int {
	for {
		if atomic.CompareAndSwapUintptr(&o.lock, 0, 1) {
			break
		}
	}

	length := len(o.i)

	atomic.StoreUintptr(&o.lock, 0)
	return length
}
