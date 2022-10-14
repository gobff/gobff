package syncmap

import "sync"

type Map[V any] struct {
	data  map[string]V
	mutex *sync.Mutex
}

func New[V any]() Map[V] {
	return Map[V]{
		data:  make(map[string]V),
		mutex: new(sync.Mutex),
	}
}

func (m Map[V]) Set(key string, v V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data[key] = v
}

func (m Map[V]) Get(key string) V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.data[key]
}

func (m Map[V]) Data() map[string]V {
	return m.data
}
