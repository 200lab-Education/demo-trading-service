package locker

import "time"

type Locker interface {
	Lock(key string, expTime time.Duration) error
	Unlock(key string) error
	RLock(key string) error
}
