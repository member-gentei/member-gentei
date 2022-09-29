package roles

import (
	"sync"
)

type DefaultMapRWMutex struct {
	mapMutex *sync.Mutex
	mutexes  map[string]*sync.RWMutex
}

func (d *DefaultMapRWMutex) GetOrCreate(key string) *sync.RWMutex {
	d.mapMutex.Lock()
	defer d.mapMutex.Unlock()
	mutex, exists := d.mutexes[key]
	if exists {
		return mutex
	}
	mutex = &sync.RWMutex{}
	d.mutexes[key] = mutex
	return mutex
}

func NewDefaultMapRWMutex() *DefaultMapRWMutex {
	return &DefaultMapRWMutex{
		mapMutex: &sync.Mutex{},
		mutexes:  map[string]*sync.RWMutex{},
	}
}
