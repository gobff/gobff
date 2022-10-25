package syncmap

import "sync"

type Map[Value any] struct {
	data  map[string]Value
	mutex *sync.Mutex
}

func New[Value any]() Map[Value] {
	return Map[Value]{
		data:  make(map[string]Value),
		mutex: new(sync.Mutex),
	}
}

func (m Map[Value]) Set(key string, v Value) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data[key] = v
}

func (m Map[Value]) Get(key string) Value {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.data[key]
}

func (m Map[Value]) Data() map[string]Value {
	return m.data
}
