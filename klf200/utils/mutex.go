package utils

import (
	"context"
)

// A Mutex is a mutual exclusion lock.
type Mutex interface {
	// Lock locks m.
	// If the lock is already in use, the calling goroutine
	// blocks until the mutex is available.
	Lock()

	// TryLock tries to lock m and reports whether it succeeded.
	//
	// Note that while correct uses of TryLock do exist, they are rare,
	// and use of TryLock is often a sign of a deeper problem
	// in a particular use of mutexes.
	TryLock() bool

	// TryLock tries to lock m and reports whether it succeeded.
	TryLockWithContext(ctx context.Context) bool

	// Unlock unlocks m.
	// It is a run-time error if m is not locked on entry to Unlock.
	//
	// A locked Mutex is not associated with a particular goroutine.
	// It is allowed for one goroutine to lock a Mutex and then
	// arrange for another goroutine to unlock it.
	Unlock()
}

func NewMutex() Mutex {
	return &mutex{
		c: make(chan struct{}, 1),
	}
}

type mutex struct {
	c chan struct{}
}

func (m *mutex) Lock() {
	m.c <- struct{}{}
}

func (m *mutex) TryLock() bool {
	select {
	case m.c <- struct{}{}:
		return true
	default:
		return false
	}
}

func (m *mutex) TryLockWithContext(ctx context.Context) bool {
	select {
	case m.c <- struct{}{}:
		return true
	case <-ctx.Done():
		return false
	}
}

func (m *mutex) Unlock() {
	<-m.c
}
