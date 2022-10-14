package donewatcher

import (
	"sync"
)

type (
	Watcher interface {
		Done(key string)
		Wait(keys []string)
	}
	watcher struct {
		keysDone  map[string]bool
		observers []*observer
		mutex     *sync.Mutex
	}
)

func NewWatcher() Watcher {
	return &watcher{
		keysDone:  make(map[string]bool),
		observers: make([]*observer, 0),
		mutex:     new(sync.Mutex),
	}
}

func (m *watcher) Done(key string) {
	m.mutex.Lock()
	m.keysDone[key] = true
	m.checkObservers()
	m.mutex.Unlock()
}

func (m *watcher) Wait(keys []string) {
	o := newObserver(keys)

	go func() {
		m.mutex.Lock()
		m.observers = append(m.observers, o)
		o.check(m.keysDone)
		m.mutex.Unlock()
	}()

	o.wait()
}

func (m *watcher) checkObservers() {
	for _, observer := range m.observers {
		observer.check(m.keysDone)
	}
}

type observer struct {
	keys    []string
	channel chan bool
	done    bool
}

func newObserver(keys []string) *observer {
	return &observer{
		keys:    keys,
		channel: make(chan bool),
	}
}

func (o *observer) check(keys map[string]bool) {
	if o.done {
		return
	}
	for _, key := range o.keys {
		if !keys[key] {
			return
		}
	}
	o.done = true
	o.channel <- true
}

func (o *observer) wait() {
	<-o.channel
	close(o.channel)
}
