package orderedmap

import "sync"

type OrderedMap[K comparable, V any] struct {
	keys   []K
	values map[K]V
	mu     sync.RWMutex
}

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func New[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys:   make([]K, 0),
		values: make(map[K]V),
	}
}

func FromMap[K comparable, V any](m map[K]V, keys []K) *OrderedMap[K, V] {
	om := New[K, V]()
	for _, k := range keys {
		if v, ok := m[k]; ok {
			om.Set(k, v)
		}
	}
	return om
}

func (o *OrderedMap[K, V]) Set(key K, value V) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if _, exists := o.values[key]; !exists {
		o.keys = append(o.keys, key)
	}
	o.values[key] = value
}

func (o *OrderedMap[K, V]) Get(key K) (value V, ok bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	v, ok := o.values[key]
	return v, ok
}

func (o *OrderedMap[K, V]) Delete(key K) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if _, exists := o.values[key]; !exists {
		return
	}
	delete(o.values, key)
	for i, k := range o.keys {
		if k == key {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
}

func (o *OrderedMap[K, V]) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.keys)
}

func (o *OrderedMap[K, V]) Keys() []K {
	o.mu.RLock()
	defer o.mu.RUnlock()

	result := make([]K, len(o.keys))
	copy(result, o.keys)
	return result
}

func (o *OrderedMap[K, V]) Index(key K) int {
	o.mu.RLock()
	defer o.mu.RUnlock()

	for i, k := range o.keys {
		if k == key {
			return i
		}
	}
	return -1
}

func (o *OrderedMap[K, V]) Values() []V {
	o.mu.RLock()
	defer o.mu.RUnlock()

	values := make([]V, 0, len(o.keys))
	for _, k := range o.keys {
		values = append(values, o.values[k])
	}
	return values
}

func (o *OrderedMap[K, V]) First() *Pair[K, V] {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if len(o.keys) == 0 {
		return nil
	}
	key := o.keys[0]
	return &Pair[K, V]{Key: key, Value: o.values[key]}
}

func (o *OrderedMap[K, V]) Last() *Pair[K, V] {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if len(o.keys) == 0 {
		return nil
	}
	key := o.keys[len(o.keys)-1]
	return &Pair[K, V]{Key: key, Value: o.values[key]}
}

func (o *OrderedMap[K, V]) Iterate(fn func(pair *Pair[K, V]) bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	for _, k := range o.keys {
		if !fn(&Pair[K, V]{Key: k, Value: o.values[k]}) {
			break
		}
	}
}
