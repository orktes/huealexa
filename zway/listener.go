package zway

import (
	"sync"
	"sync/atomic"
)

type listener struct {
	sync.RWMutex
	nextID    int64
	listeners map[int64]func(interface{})
}

func (l *listener) addListener(fn func(interface{})) int64 {
	l.Lock()
	defer l.Unlock()

	if l.listeners == nil {
		l.listeners = map[int64]func(interface{}){}
	}

	id := atomic.AddInt64(&l.nextID, 1)
	l.listeners[id] = fn
	return id
}

func (l *listener) removeListener(id int64) {
	l.Lock()
	defer l.Unlock()

	if l.listeners == nil {
		return
	}

	delete(l.listeners, id)
}

func (l *listener) emit(val interface{}) {
	l.Lock()
	defer l.Unlock()

	if l.listeners == nil {
		return
	}

	for _, fn := range l.listeners {
		fn(val)
	}
}
