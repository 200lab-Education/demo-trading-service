package local

import (
	"sync"
	"time"
)

type localLocker struct {
	lock   *sync.RWMutex
	store  map[string]*sync.RWMutex
	prefix string
}

func NewLocalLocker(prefix string) *localLocker {
	return &localLocker{
		lock:   new(sync.RWMutex),
		store:  make(map[string]*sync.RWMutex),
		prefix: prefix,
	}
}

func (l *localLocker) Lock(key string, expTime time.Duration) error {
	l.lock.Lock()

	//go func() {
	//	time.Sleep(expTime)
	//	l.Unlock(key)
	//}()

	if k, ok := l.store[key]; ok {
		l.lock.Unlock()
		k.Lock()
		return nil
	}

	newLock := new(sync.RWMutex)
	l.store[key] = newLock
	l.lock.Unlock()

	l.store[key].Lock()

	return nil
}

func (l *localLocker) Unlock(key string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if k, ok := l.store[key]; ok {
		k.Unlock()
	}

	return nil
}

func (localLocker) RLock(key string) error {
	return nil
}

func (l *localLocker) GetPrefix() string {
	return l.prefix
}

func (l *localLocker) Get() interface{} {
	return l
}

func (l *localLocker) Name() string {
	return l.GetPrefix()
}

func (l *localLocker) InitFlags() {}

func (l *localLocker) Configure() error {
	return nil
}

func (l *localLocker) Run() error {
	return l.Configure()
}

func (l *localLocker) Stop() <-chan bool {
	c := make(chan bool)
	go func() { c <- true }()
	return c
}
