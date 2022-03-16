package roles

import "sync"

type DefaultMapRWMutex struct {
	mutexes map[string]*sync.RWMutex
}

func (d *DefaultMapRWMutex) GetOrCreate(key string) *sync.RWMutex {
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
		mutexes: map[string]*sync.RWMutex{},
	}
}
